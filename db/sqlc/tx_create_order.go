package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateOrderTxParams struct {
	CreateOrderParams
	ListUserCartItemsParams
	Country string
	Address string
	Town string
	PostalCode string
	Landmark string
}

type CreateOrderTxResult struct {
	Order Order
	OrderItem OrderItem
	ShippingAddress ShippingAddress
	User User
	CartArray []Cart
	Cart Cart
	Product Product
	UpdatedProduct Product
}

func (store *SQLStore) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (CreateOrderTxResult, error) {
	var result CreateOrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the Shipping Address
		result.ShippingAddress, err = q.CreateShippingAddress(ctx, CreateShippingAddressParams{
			Username: arg.Username,
			Country: arg.Country,
			Address: arg.Address,
			Town: arg.Town,
			PostalCode: pgtype.Text{
				String: arg.PostalCode,
				Valid: true,
			},
			Landmark: pgtype.Text{
				String: arg.Landmark,
				Valid: true,
			},
		})

		if err != nil {
			return err
		}

		// Create the Order
		result.Order, err = q.CreateOrder(ctx, CreateOrderParams{
			RefNo: arg.RefNo,
			Username: arg.CreateOrderParams.Username,
			Amount: arg.Amount,
			PaymentMethod: arg.PaymentMethod,
			OrderStatus: arg.OrderStatus,
			ShippingMethod: arg.ShippingMethod,
			ShippingAddressID: result.ShippingAddress.ID,
		})

		if err != nil {
			return err
		}

		// Fetch the user cart items
		result.CartArray, err = q.ListUserCartItems(ctx, ListUserCartItemsParams{
			UserRefID: arg.UserRefID,
			Limit: arg.Limit,
			Offset: arg.Offset,
		})

		if err != nil {
			return err
		}

		for _, item := range result.CartArray {
			// Fetch the particular product
			result.Product, err = q.GetProduct(ctx, int64(item.ProductID))

			if err != nil {
				return err
			}

			// Create the order items
			result.OrderItem, err = q.CreateOrderItem(ctx, CreateOrderItemParams{
				OrderID: result.Order.RefNo,
				ProductID: result.Product.ProductRefNo,
				Quantity: int32(item.ProductQuantity),
				Price: item.ProductPrice,
			})

			if err != nil {
				return err
			}

			// delete the user cart
			err = q.DeleteCartItem(ctx, item.ProductID)
			if err != nil {
				return err
			}

			// Fetch the product
			result.Product, err = q.GetProduct(ctx, int64(item.ProductID))

			if err != nil {
				return err
			}

			currentProductCount := result.Product.Quantity

			// Update the new product count
			result.UpdatedProduct, err = q.UpdateProduct(ctx, UpdateProductParams{
				ID: int64(item.ProductID),
				Quantity: pgtype.Int4{
					Int32: currentProductCount - int32(item.ProductQuantity), 
					Valid: true,
				},
			})

			if err != nil {
				return err
			}
		}

		return err
	})
	return result, err
}