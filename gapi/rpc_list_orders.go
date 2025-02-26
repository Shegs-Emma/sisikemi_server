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

func (server *Server) ListOrders(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if !authPayload.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	arg := db.ListOrdersParams{
		Limit: req.GetPageSize(),
		Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListOrders(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	var pbOrders []*pb.Order
	for _, item := range result {
		pbOrders = append(pbOrders, &pb.Order{
			Id: item.ID,
			RefNo: item.RefNo,
			Username: fetchReferencedUserByUsername(server, ctx, item.Username),
			Amount: float64(item.Amount),
			PaymentMethod: item.PaymentMethod,
			ShippingAddressId: fetchReferencedShippingAddress(server, ctx, item.ShippingAddressID),
			ShippingMethod: item.ShippingMethod,
			OrderStatus: string(item.OrderStatus),
			Items: fetchReferencedOrderItems(server, ctx, item.RefNo),
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}
	
	rsp := &pb.ListOrderResponse{
		Orders: pbOrders,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
	}

	return rsp, nil
}