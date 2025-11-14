package gapi

import (
	"context"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateProduct (ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if !authPayload.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to make this action")
	}

	arg := db.CreateProductTxParams{
		CreateProductParams: db.CreateProductParams{
			ProductRefNo: util.RandomString(12),
			ProductName: req.GetProductName(),
			ProductDescription: req.GetProductDescription(),
			ProductCode: req.GetProductCode(),
			Price: req.GetPrice(),
			SalePrice: req.GetSalePrice(),
			Collection: req.GetCollection(),
			Quantity: req.GetQuantity(),
			Color: req.GetColor(),
			Size: req.GetSize(),
			Status: db.ProductStatus(req.GetStatus()),
		},
		MainImage: req.GetMainImage(),
		OtherImage1: req.GetOtherImage_1(),
		OtherImage2: req.GetOtherImage_2(),
		OtherImage3: req.GetOtherImage_3(),
	}

	txResult, err := server.store.CreateProductTx(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create product: %s", err)
	}

	rsp := &pb.CreateProductResponse{
		Product: convertProduct(server, ctx, txResult.Product),
	}

	return rsp, nil
}