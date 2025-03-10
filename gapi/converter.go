package gapi

import (
	"context"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func convertUser(user db.User) *pb.User {
	return &pb.User{
		Id: user.ID,
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

func convertShippingAddress(server *Server, ctx context.Context, addr db.ShippingAddress) *pb.ShippingAddress {
	return &pb.ShippingAddress{
		Username: fetchReferencedUserByUsername(server, ctx, addr.Username),
		Country: addr.Country,
		Address: addr.Address,
		Town: addr.Town,
		PostalCode: addr.PostalCode.String,
		Landmark: addr.Landmark.String,
		CreatedAt: timestamppb.New(addr.CreatedAt),
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

func convertOrder(server *Server, ctx context.Context, order db.Order) *pb.Order {
	return &pb.Order{
		Id: order.ID,
		RefNo: order.RefNo,
		Username: fetchReferencedUserByUsername(server, ctx, order.Username),
		Amount: float64(order.Amount),
		PaymentMethod: order.PaymentMethod,
		ShippingAddressId: fetchReferencedShippingAddress(server, ctx, order.ShippingAddressID),
		ShippingMethod: order.ShippingMethod,
		OrderStatus: string(order.OrderStatus),
		Items: fetchReferencedOrderItems(server, ctx, order.RefNo),
		CreatedAt: timestamppb.New(order.CreatedAt),
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

func fetchReferencedShippingAddress(server *Server, ctx context.Context,  addr int64) *pb.ShippingAddress {
	referencedShippingAddress, err := server.store.GetShippingAddress(ctx, addr)

	if err != nil {
		return nil
	}

	return convertShippingAddress(server, ctx, referencedShippingAddress)
}

func fetchReferencedUserByUsername(server *Server, ctx context.Context,  user string) *pb.User {
	referencedUser, err := server.store.GetUserByUsername(ctx, user)

	if err != nil {
		return nil
	}

	return convertUser(referencedUser)
}

func fetchReferencedOrderItems(server *Server, ctx context.Context,  orderRef string) []*pb.OrderItem {
	orderItems, err := server.store.GetOrderItemsForOrder(ctx, db.GetOrderItemsForOrderParams{
		OrderID: orderRef,
		Limit: 10,
	})

	if err != nil {
		return nil
	}

	var pbUserOrderItems []*pb.OrderItem
	for _, item := range orderItems {
		pbUserOrderItems = append(pbUserOrderItems, &pb.OrderItem{
			Id: item.ID,
			OrderId: item.OrderID,
			ProductId: item.ProductID,
			Quantity: int64(item.Quantity),
			Price: float64(item.Price),
			CreatedAt: timestamppb.New(item.CreatedAt),
		})
	}

	return pbUserOrderItems;
}

func parseSizeArray(size []string) []string {
	return size
}

func parseColorArray(color []string) []string {
	return color
}
