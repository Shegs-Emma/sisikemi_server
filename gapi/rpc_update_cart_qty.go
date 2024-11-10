package gapi

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
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

	// Fetch the current user id
	fetchedUser, err := server.store.GetUser(ctx, authPayload.Username);

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