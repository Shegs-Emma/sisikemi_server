package gapi

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) ListProducts (ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload == nil {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to update this user")
	}

	arg := db.ListProductsParams{
		Limit: req.GetPageSize(),
		Offset:  (req.GetPageId() - 1) * req.GetPageSize(),
	}

	result, err := server.store.ListProducts(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
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
			Collection: item.Collection,
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

	rsp := &pb.ListProductResponse{
		Product: pbProductItems,
		NextPageToken: fmt.Sprintf("%d", req.GetPageId() + 1),
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

func textToString(text pgtype.Text) string {
    if text.Valid {
        return text.String
    }
    return ""
}