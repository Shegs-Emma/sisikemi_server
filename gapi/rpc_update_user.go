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

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if !authPayload.IsAdmin && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FirstName: pgtype.Text{
			String: req.GetFirstName(),
			Valid: req.FirstName != nil,
		},
		LastName: pgtype.Text{
			String: req.GetLastName(),
			Valid: req.LastName != nil,
		},
		PhoneNumber: pgtype.Text{
			String: req.GetPhoneNumber(),
			Valid: req.PhoneNumber != nil,
		},
		ProfilePhoto: pgtype.Text{
			String: req.GetProfilePhoto(),
			Valid: req.ProfilePhoto != nil,
		},
	}

	if req.Password != nil {
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

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
	}
	
	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if req.FirstName != nil {
		if err := val.ValidateFirstName(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation("first_name", err))
		}
	}	

	if req.LastName != nil {
		if err := val.ValidateLastName(req.GetLastName()); err != nil {
			violations = append(violations, fieldViolation("last_name", err))
		}
	}	
	
	return violations
}