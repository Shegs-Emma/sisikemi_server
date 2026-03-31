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
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 5 
	}

	pageID := req.GetPageId()
	if pageID <= 0 {
		pageID = 1
	}

	arg := db.ListProductsParams{
		Limit:       pageSize,
		Offset:      (pageID - 1) * pageSize,
		Search:      ToPgText(req.GetSearch()),
		Collection:  ToPgInt8(req.GetCollection()),
		MinPrice:    ToPgInt8(req.GetMinPrice()),
		MaxPrice:    ToPgInt8(req.GetMaxPrice()),
		ProductName: ToPgText(req.GetProductName()),
		SortField:   ToPgText(req.GetSortBy()), 
		SortOrder:   ToPgText(req.GetSortDir()),
	}

	// Fetch Data (Single Query with Joins)
	result, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %s", err.Error())
	}

	var pbProductItems []*pb.Product

	// Map Data (No DB calls inside this loop!)
	for _, item := range result {
		
		var collectionObj *pb.Collection
		if item.CollectionName.Valid {
			collectionObj = &pb.Collection{
				Id:   item.Collection, 
				CollectionName: item.CollectionName.String,
			}
		}

		pbProductItems = append(pbProductItems, &pb.Product{
			Id:                  item.ID,
			ProductRefNo:        item.ProductRefNo,
			ProductName:         item.ProductName,
			ProductDescription:  item.ProductDescription,
			ProductCode:         item.ProductCode,
			Price:               item.Price,
			SalePrice:           item.SalePrice,
			Collection:          collectionObj, // Populated from JOIN
			Quantity:            item.Quantity,
			Color:               item.Color,
			Size:                item.Size,
			Status:              string(item.Status),
			ProductImageMain: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageMain)),
			ProductImageOther_1: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther1)),
			ProductImageOther_2: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther2)),
			ProductImageOther_3: fetchReferencedProductMedium(server, ctx, textToString(item.ProductImageOther3)),
			CreatedAt:           timestamppb.New(item.CreatedAt),
		})
	}

	// Get Total Count
	totalCount, err := server.store.CountProducts(ctx, db.CountProductsParams{
		Search:      ToPgText(req.GetSearch()),
		Collection:  ToPgInt8(req.GetCollection()),
		MinPrice:    ToPgInt8(req.GetMinPrice()),
		MaxPrice:    ToPgInt8(req.GetMaxPrice()),
		ProductName: ToPgText(req.GetProductName()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "count products error: %s", err)
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(pageSize)))

	var nextPageToken string
	if pageID < totalPages {
		nextPageToken = fmt.Sprintf("%d", pageID+1)
	}

	rsp := &pb.ListProductResponse{
		Product:       pbProductItems,
		NextPageToken: nextPageToken,
		TotalPages:    totalPages,
		TotalCount:    int32(totalCount),
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

func ToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func ToPgInt(i int32) pgtype.Int4 {
	if i == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: i, Valid: true}
}

func ToPgInt8(i int64) pgtype.Int8 {
    if i == 0 { 
        return pgtype.Int8{Valid: false}
    }
    return pgtype.Int8{Int64: i, Valid: true} // Int64 matches BIGINT
}