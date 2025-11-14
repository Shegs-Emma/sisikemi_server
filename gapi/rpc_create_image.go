package gapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	// Create a path to save the uploaded image
	imagePath := filepath.Join("uploads", req.Filename)
	
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(imagePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Write the image to a file
	if err := os.WriteFile(imagePath, req.Image, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image: %v", err)
	}

	arg := db.CreateMediaParams{
		MediaRef: util.RandomString(10),
		Url: fmt.Sprintf("http://localhost:8080/uploads/%s", req.Filename),
		AwsID: util.RandomString(12),
	}

	result, err := server.store.CreateMedia(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.UploadImageResponse{
		Media: convertMedia(result),
	}

	return rsp, nil
}