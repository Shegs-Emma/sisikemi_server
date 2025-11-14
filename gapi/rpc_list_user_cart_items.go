package gapi

import (
	"context"
	"fmt"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListUserCartItems (ctx context.Context, req *pb.ListUserCartRequest) (*pb.ListUserCartResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	// fetch the user id
	fetchedUser, err := server.store.GetUserByUsername(ctx, authPayload.Username);

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user could not be fetched")
	}

	arg := db.ListUserCartItemsParams{
		UserRefID: fetchedUser.ID,
		Limit: req.GetPageSize(),
		Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListUserCartItems(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	var pbUserCartItems []*pb.Cart
	for _, item := range result {
		pbUserCartItems = append(pbUserCartItems, &pb.Cart{
			Id: item.ID,
			ProductId: int64(item.ProductID),
			ProductName: item.ProductName,
			ProductPrice: item.ProductPrice,
			ProductQuantity: item.ProductQuantity,
			ProductImage: item.ProductImage,
			ProductColor: item.ProductColor,
			ProductSize: item.ProductSize,
			UserRefId: fetchReferencedUser(server, ctx, item.UserRefID),
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	rsp := &pb.ListUserCartResponse{
		Cart: pbUserCartItems,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
	}

	return rsp, nil
}