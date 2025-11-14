package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	_ "github.com/Shegs-Emma/sisikemi_server/doc/statik"
	"github.com/Shegs-Emma/sisikemi_server/gapi"
	"github.com/Shegs-Emma/sisikemi_server/mail"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/Shegs-Emma/sisikemi_server/worker"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Msg("Cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	cloud := util.SetupCloudinary(config)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor, cloud)

	runGrpcServer(ctx, waitGroup, config, store, taskDistributor, cloud)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	opts := mail.Mail{
		Domain: config.MailDomain,
		Host: config.MailHost,
		Port: config.MailPort,
		Username: config.MailUsername,
		Password: config.MailPassword,
		Encryption: config.MailEncryption,
		FromAddress: config.FromAddress,
		FromName: config.FromName,
	}

	smtpMailer := mail.NewSendSMTPEmailSender(opts)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, smtpMailer, config)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

// UploadImageHandler is the handler for file uploads with pathParams
func uploadImageHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	fmt.Println("err", err)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the posted form-data
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file's content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Get the filename from form data or pathParams if provided
	filename := r.FormValue("filename")
	if filename == "" {
		filename = handler.Filename // Default to uploaded filename if not provided
	}

	// Use pathParams if needed (for example, adding to the filename)
	if val, ok := pathParams["id"]; ok {
		filename = val + "_" + filename // Modify the filename with pathParam
	}

	// Create the file path to save
	imagePath := filepath.Join("uploads", filename)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(imagePath), 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Write the image to a file
	err = os.WriteFile(imagePath, fileBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Respond with the URL of the uploaded image
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Image uploaded successfully", "url": "http://localhost:8080/uploads/%s"}`, filename)
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor, cloud *cloudinary.Cloudinary) {
	server, err := gapi.NewServer(config, store, taskDistributor, cloud)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSisikemiFashionServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)

		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor, cloud *cloudinary.Cloudinary) {
	server, err := gapi.NewServer(config, store, taskDistributor, cloud)

	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	// Use the wrapper for the upload image handler
	grpcMux.HandlePath("POST", "/v1/upload_image", uploadImageHandler)

	err = pb.RegisterSisikemiFashionHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	// Serve static files from the uploads directory
	uploadsHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/")))
	mux.Handle("/uploads/", uploadsHandler) // Add this line to serve uploads


	c := cors.New(cors.Options{
		AllowedOrigins: config.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	})

	handler := c.Handler(gapi.HttpLogger(mux))

	httpServer := &http.Server{
		Handler: handler,
		Addr: config.HttpServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()

		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}

		log.Info().Msg("Http gateway server is stopped")
		return nil
	})
}