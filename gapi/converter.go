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
		ProfilePhoto: user.ProfilePhoto.String,
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
		ProductCount: collection.ProductCount.Int64,
		CreatedAt: timestamppb.New(collection.CreatedAt),
	}
}

func convertProduct(server *Server, ctx context.Context, product db.Product) *pb.Product {
	return &pb.Product{
		Id: product.ID,
		ProductRefNo: product.ProductRefNo,
		ProductName: product.ProductName,
		ProductDescription: product.ProductDescription,
		ProductCode: product.ProductCode,
		Price: product.Price,
		SalePrice: product.SalePrice,
		Collection: fetchReferencedCollection(server, ctx, product.Collection),
		Quantity: product.Quantity,
		Color: parseColorArray(product.Color),
		Size: parseSizeArray(product.Size),
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

func convertCart(server *Server, ctx context.Context, cart db.Cart) *pb.Cart {
	return &pb.Cart{
		Id: cart.ID,
		ProductId: int64(cart.ProductID),
		ProductName: cart.ProductName,
		ProductPrice: cart.ProductPrice,
		ProductQuantity: cart.ProductQuantity,
		ProductImage: cart.ProductImage,
		ProductColor: cart.ProductColor,
		ProductSize: cart.ProductSize,
		UserRefId: fetchReferencedUser(server, ctx, cart.UserRefID),
		CreatedAt: timestamppb.New(cart.CreatedAt),
	}
}

func fetchReferencedMedia(server *Server, ctx context.Context,  media string) *pb.Media {
	referencedMedia, err := server.store.GetMediaByRef(ctx, media)

	if err != nil {
		return nil
	}

	return convertMedia(referencedMedia)
}

func fetchReferencedUser(server *Server, ctx context.Context,  user int64) *pb.User {
	referencedUser, err := server.store.GetUserById(ctx, user)

	if err != nil {
		return nil
	}

	return convertUser(referencedUser)
}

func parseSizeArray(size []string) []string {
	return size
}

func parseColorArray(color []string) []string {
	return color
}