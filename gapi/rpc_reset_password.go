package gapi

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	violations := validateVerifyResetRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	result, err := server.store.GetUser(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}

	// Check if the code is the same
	if req.GetCode() != result.VerificationCode.String {
		return nil, status.Errorf(codes.PermissionDenied, "%s", err.Error())
	}

	arg := db.UpdateUserParams{
		Username: result.Username,
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid: true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time: time.Now(),
			Valid: true,
		}
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
	}

	rsp := &pb.ResetPasswordResponse{
		IsUpdated: true,
	}

	return rsp, nil
}

func validateVerifyResetRequest(req *pb.ResetPasswordRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateCode(req.GetCode()); err != nil {
		violations = append(violations, fieldViolation("code", err))
	}

	return violations
}