package gapi

import (
	"context"
	"time"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/Shegs-Emma/sisikemi_server/val"
	"github.com/Shegs-Emma/sisikemi_server/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	violations := validateForgotPasswordRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	verificationCode := util.RandomVerificationCode(6)

	arg := db.ForgotPasswordTxParams{
		UpdateUserVerificationCodeParams: db.UpdateUserVerificationCodeParams{
			Email: req.GetEmail(),
			VerificationCode: pgtype.Text{
				String: verificationCode,
				Valid: true,
			},
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerificationCodeEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskdistributor.DistributeTaskSendVerificationCodeEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.ForgotPasswordTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.ForgotPasswordResponse{
		CodeIsSent: txResult.User.VerificationCode.Valid,
	}

	return rsp, nil
}

func validateForgotPasswordRequest(req *pb.ForgotPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}