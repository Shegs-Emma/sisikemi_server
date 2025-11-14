package gapi

import (
	"context"

	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetProduct (ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	result, err := server.store.GetProduct(ctx, req.GetId())

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	rsp := &pb.GetProductResponse{
		Product: convertProduct(server, ctx, result),
	}

	return rsp, nil
}