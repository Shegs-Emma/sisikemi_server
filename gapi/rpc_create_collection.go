package gapi

import (
	"context"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateCollection (ctx context.Context ,req *pb.CreateCollectionRequest) (*pb.CreateCollectionResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if !authPayload.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	arg := db.CreateCollectionParams{
		CollectionName: req.GetCollectionName(),
		CollectionDescription: req.GetCollectionDescription(),
		ThumbnailImage: req.GetThumbnailImage(),
		HeaderImage: req.GetHeaderImage(),
	}

	result, err := server.store.CreateCollection(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateCollectionResponse{
		Collection: convertCollection(result),
	}
	
	return rsp, nil
}