package gapi

import (
	"context"
	"fmt"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListMedia (ctx context.Context, req *pb.ListMediaRequest) (*pb.ListMediaResponse, error) {
	arg := db.ListMediaParams{
		Limit: req.GetPageSize(),
		Offset: (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListMedia(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	var pbMediaItems []*pb.Media
	for _, item := range result {
		pbMediaItems = append(pbMediaItems, &pb.Media{
			Id: item.ID,
			MediaRef: item.MediaRef,
			Url: item.Url,
			AwsId: item.AwsID,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	rsp := &pb.ListMediaResponse{
		Media: pbMediaItems,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
	}

	return rsp, nil
}