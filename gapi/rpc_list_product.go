package gapi

import (
	"context"
	"fmt"
	"math"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListProducts (ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	arg := db.ListProductsParams{
		Limit: req.GetPageSize(),
		Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListProducts(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}

	var pbProductItems []*pb.Product
	for _, item := range result {
		pbProductItems = append(pbProductItems, &pb.Product{
			Id: item.ID,
			ProductRefNo: item.ProductRefNo,
			ProductName: item.ProductName,
			ProductDescription: item.ProductDescription,
			ProductCode: item.ProductCode,
			Price: item.Price,
			SalePrice: item.SalePrice,
			Collection: fetchReferencedCollection(server, ctx, item.Collection),
			Quantity: item.Quantity,
			Color: item.Color,
			Size: item.Size,
			Status: string(item.Status),
			ProductImageMain: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageMain)),
			ProductImageOther_1: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther1)),
			ProductImageOther_2: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther2)),
			ProductImageOther_3: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther3)),
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	totalCount, err := server.store.CountProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "count products error: %s", err)
	}

	rsp := &pb.ListProductResponse{
		Product: pbProductItems,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
		TotalPages: int32(math.Ceil(float64(totalCount) / float64(req.PageSize))),
	}

	return rsp, nil
}

func fetchReferencedProductMedium(server *Server, ctx context.Context,  media string) *pb.ProductMedium {
	productMedia, err := server.store.GetProductMediaByRef(ctx, media)

	if err != nil {
		return nil
	}

	return convertProductMedium(server, ctx, productMedia)
}

func fetchReferencedCollection(server *Server, ctx context.Context,  collection int64) *pb.Collection {
	productCollection, err := server.store.GetCollection(ctx, collection)

	if err != nil {
		return nil
	}

	return convertCollection(productCollection)
}

func textToString(text pgtype.Text) string {
    if text.Valid {
        return text.String
    }
    return ""
}