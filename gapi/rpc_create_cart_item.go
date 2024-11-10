package gapi

import (
	"context"
	"fmt"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateCartItem (ctx context.Context, req *pb.CreateCartItemRequest) (*pb.CreateCartItemResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	fmt.Print(authPayload)

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	// Check if there is already a cart item with that product
	_, err = server.store.GetCartItemByProductId(ctx, int32(req.GetProductId()))

	if err == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you already added this product to cart")
	}

	// Fetch the current user info
	fetchedUser, err := server.store.GetUser(ctx, authPayload.Username);

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user does not exist.")
	}

	arg := db.CreateCartItemParams{
		ProductID: int32(req.GetProductId()),
		ProductName: req.GetProductName(),
		ProductPrice: req.GetProductPrice(),
		ProductQuantity: req.GetProductQuantity(),
		ProductImage: req.GetProductImage(),
		ProductColor: req.GetProductColor(),
		ProductSize: req.GetProductSize(),
		UserRefID: fetchedUser.ID,
	}

	result, err := server.store.CreateCartItem(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateCartItemResponse{
		Cart: convertCart(server, ctx, result),
	}

	return rsp, nil
}