package gapi

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateOrder (ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to create an order")
	}

	arg := db.CreateOrderTxParams{
		CreateOrderParams: db.CreateOrderParams{
			RefNo: util.RandomString(8),
			Username: req.GetUsername(),
			Amount: req.GetAmount(),
			PaymentMethod: req.GetPaymentMethod(),
			OrderStatus: db.OrderStatus(req.GetOrderStatus()),
			ShippingMethod: req.GetShippingMethod(),
		},
		ListUserCartItemsParams: db.ListUserCartItemsParams{
			UserRefID: req.GetUserRefId(),
			Limit: req.GetPageSize(),
			Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
		},
		Country: req.GetCountry(),
		Address: req.GetAddress(),
		Town: req.GetTown(),
		PostalCode: req.GetPostalCode(),
		Landmark: req.GetLandmark(),
	}

	txResult, err := server.store.CreateOrderTx(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create product: %s", err)
	}

	rsp := &pb.CreateOrderResponse{
		Order: convertOrder(server, ctx, txResult.Order),
	}

	return rsp, nil
}