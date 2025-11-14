package gapi

import (
	"context"

	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetMedia (ctx context.Context, req *pb.GetMediaRequest) (*pb.GetMediaResponse, error) {
	result, err := server.store.GetMedia(ctx, req.GetId())

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}

	rsp := &pb.GetMediaResponse{
		Media: convertMedia(result),
	}

	return rsp, nil
}