package gapi

import (
	"context"

	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.VerifyCodeResponse, error) {
	violations := validateVerifyCodeRequest(req)
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

	rsp := &pb.VerifyCodeResponse{
		IsVerified: true,
	}

	return rsp, nil
}

func validateVerifyCodeRequest(req *pb.VerifyCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateCode(req.GetCode()); err != nil {
		violations = append(violations, fieldViolation("code", err))
	}

	return violations
}