package gapi

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		ProfilePhoto: user.ProfilePhoto,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
		IsAdmin: user.IsAdmin,
	}
}

func convertMedia(media db.Medium) *pb.Media {
	return &pb.Media{
		Id: media.ID,
		MediaRef: media.MediaRef,
		Url: media.Url,
		AwsId: media.AwsID,
		CreatedAt: timestamppb.New(media.CreatedAt),
	}
}

func convertCollection(collection db.Collection) *pb.Collection {
	return &pb.Collection{
		CollectionName: collection.CollectionName,
		CollectionDescription: collection.CollectionDescription,
		ThumbnailImage: collection.ThumbnailImage,
		HeaderImage: collection.HeaderImage,
		CreatedAt: timestamppb.New(collection.CreatedAt),
	}
}

func convertProduct(server *Server, ctx context.Context, product db.Product) *pb.Product {
	return &pb.Product{
		ProductRefNo: product.ProductRefNo,
		ProductName: product.ProductName,
		ProductDescription: product.ProductDescription,
		ProductCode: product.ProductCode,
		Price: product.Price,
		SalePrice: product.SalePrice,
		Collection: product.Collection,
		Quantity: product.Quantity,
		Color: product.Color,
		Size: product.Size,
		Status: string(product.Status),
		ProductImageMain: fetchReferencedProductMedium(server, ctx, textToString(product.ProductImageMain)),
		ProductImageOther_1: fetchReferencedProductMedium(server, ctx, textToString(product.ProductImageOther1)),
		ProductImageOther_2: fetchReferencedProductMedium(server, ctx, textToString(product.ProductImageOther2)),
		ProductImageOther_3: fetchReferencedProductMedium(server, ctx, textToString(product.ProductImageOther3)),
		CreatedAt: timestamppb.New(product.CreatedAt),
	}
}

// Convert db.ProductMedium to pb.ProductMedium
func convertProductMedium(server *Server, ctx context.Context, media db.ProductMedium) *pb.ProductMedium {
    return &pb.ProductMedium{
        Id: media.ID,
		ProductMediaRef: media.ProductMediaRef,
		ProductId: media.ProductID,
		IsMainImage: media.IsMainImage,
		MediaId: fetchReferencedMedia(server, ctx, media.MediaID),
    }
}

func fetchReferencedMedia(server *Server, ctx context.Context,  media string) *pb.Media {
	referencedMedia, err := server.store.GetMediaByRef(ctx, media)

	if err != nil {
		return nil
	}

	return convertMedia(referencedMedia)
}