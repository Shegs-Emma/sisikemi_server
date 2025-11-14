package gapi

import (
	"context"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMedia (ctx context.Context, req *pb.CreateMediaRequest) (*pb.CreateMediaResponse, error) {
	arg := db.CreateMediaParams{
		MediaRef: req.GetMediaRef(),
		Url: req.GetUrl(),
		AwsID: req.GetAwsId(),
	}

	result, err := server.store.CreateMedia(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateMediaResponse{
		Media: convertMedia(result),
	}

	return rsp, nil
}