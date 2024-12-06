package gapi

import (
	"context"
	"fmt"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteProduct (ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if !authPayload.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to make this action")
	}

	arg := db.DeleteProductTxParams{
		ProductId: req.GetProductId(),
	}

	_, err = server.store.DeleteProductTx(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %s", err)
	}

	rsp := &pb.DeleteProductResponse{
		Message: fmt.Sprintln("Product has been deleted successfully"),
	}

	return rsp, nil
}