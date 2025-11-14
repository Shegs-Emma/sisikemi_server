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

func (server *Server) UpdateProduct(ctx context.Context,req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if !authPayload.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update product")
	}

	arg := db.UpdateProductParams{
		ID: req.GetProductId(),
		ProductName: pgtype.Text{
			String: req.GetProductName(),
			Valid: req.ProductName != nil,
		},
		ProductDescription: pgtype.Text{
			String: req.GetProductDescription(),
			Valid: req.ProductDescription != nil,
		},
		ProductCode: pgtype.Text{
			String: req.GetProductCode(),
			Valid: req.ProductCode != nil,
		},
		Price: pgtype.Int8{
			Int64: req.GetPrice(),
			Valid: req.Price != nil,
		},
		SalePrice: pgtype.Text{
			String: req.GetSalePrice(),
			Valid: req.SalePrice != nil,
		},
		Collection: pgtype.Int8{
			Int64: req.GetCollection(),
			Valid: req.Collection != nil,
		},
		Quantity: pgtype.Int4{
			Int32: req.GetQuantity(),
			Valid: req.Quantity != nil,
		},
		Color: req.GetColor(),
		Size: req.GetSize(),
		Status: db.NullProductStatus{
			ProductStatus: db.ProductStatus(req.GetStatus()),
			Valid: req.Status != nil,
		},
		// ProductImageMain: pgtype.Text{
		// 	String: req.GetMainImage(),
		// 	Valid: true,
		// },
		// ProductImageOther1: pgtype.Text{
		// 	String: req.GetOtherImage_1(),
		// 	Valid: true,
		// },
		// ProductImageOther2: pgtype.Text{
		// 	String: req.GetOtherImage_2(),
		// 	Valid: true,
		// },
		// ProductImageOther3: pgtype.Text{
		// 	String: req.GetOtherImage_3(),
		// 	Valid: true,
		// },
	}

	updatedProduct, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to Update product: %s", err)
	}

	rsp := &pb.UpdateProductResponse{
		Product: convertProduct(server, ctx, updatedProduct),
	}

	return rsp, nil
}