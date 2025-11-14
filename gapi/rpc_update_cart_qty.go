package gapi

import (
	"context"
	"errors"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateCartItemQty (ctx context.Context, req *pb.UpdateCartItemQtyRequest) (*pb.UpdateCartItemQtyResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	// Fetch the product in question to check the available quantity
	fetchedProduct, err := server.store.GetProduct(ctx, req.GetProductId())

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "product could not be fetched")
	}

	// Fetch the current user id
	fetchedUser, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user could not be fetched")
	}

	// Fetch the cart item to be updated
	fetchedCart, err := server.store.GetCartItem(ctx, req.GetItemId())

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cart item does not exist")
	}

	// Check that the current user created the cart item
	if fetchedUser.ID != fetchedCart.UserRefID {
		return nil, status.Errorf(codes.PermissionDenied, "you can't update this cart item")
	}

	// Check that there is enough products
	if int64(fetchedProduct.Quantity) < req.GetProductQuantity() && req.GetAction() == "increment" {
		return nil, status.Errorf(codes.PermissionDenied, "you dont have enough quantity")
	}

	arg := db.UpdateCartItemQtyParams{
		ID: req.GetItemId(),
		ProductQuantity: pgtype.Int8{
			Int64: req.GetProductQuantity(),
			Valid: req.ProductQuantity != nil,
		},
	}

	result, err := server.store.UpdateCartItemQty(ctx, arg)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
	}

	rsp := &pb.UpdateCartItemQtyResponse{
		Cart: convertCart(server, ctx, result),
	}

	return rsp, nil
}