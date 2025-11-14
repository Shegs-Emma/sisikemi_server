package gapi

import (
	"context"
	"fmt"

	"github.com/Shegs-Emma/sisikemi_server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteCartItem (ctx context.Context, req *pb.DeleteCartItemRequest) (*pb.DeleteCartItemResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	// Fetch the current user info
	fetchedUser, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user could not be fetched")
	}

	// Fetch the cart item to be deleted
	fetchedCart, err := server.store.GetCartItemByProductId(ctx, req.GetProductId())

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cart item does not exist")
	}

	// Check that the current user created the cart item
	if fetchedUser.ID != fetchedCart.UserRefID {
		return nil, status.Errorf(codes.PermissionDenied, "you can't delete this cart item")
	}

	err = server.store.DeleteCartItem(ctx, req.GetProductId())

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cart item could not be deleted")
	}

	rsp := &pb.DeleteCartItemResponse{
		Message: fmt.Sprintf("%s has been deleted successfully", fetchedCart.ProductName),
	}

	return rsp, nil
}