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

func (server *Server) ListCollections (ctx context.Context, req *pb.ListCollectionRequest) (*pb.ListCollectionResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	arg := db.ListCollectionParams{
		Limit: req.GetPageSize(),
		Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListCollection(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	var pbCollectionItems []*pb.Collection
	for _, item := range result {
		pbCollectionItems = append(pbCollectionItems, &pb.Collection{
			Id: item.ID,
			CollectionName: item.CollectionName,
			ProductCount: item.ProductCount.Int64,
			CollectionDescription: item.CollectionDescription,
			ThumbnailImage: item.ThumbnailImage,
			HeaderImage: item.HeaderImage,
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	rsp := &pb.ListCollectionResponse{
		Collection: pbCollectionItems,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
	}

	return rsp, nil
}