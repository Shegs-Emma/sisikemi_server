package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/techschool/simplebank/util"
)

type CreateProductTxParams struct {
	CreateProductParams
	MainImage string
    OtherImage1 string
    OtherImage2 string
    OtherImage3 string
}

type CreateProductTxResult struct {
	Product Product
	ProductMedium ProductMedium
	Collection Collection
}

func (store *SQLStore) CreateProductTx (ctx context.Context, arg CreateProductTxParams) (CreateProductTxResult, error) {
	var result CreateProductTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Product, err = q.CreateProduct(ctx, CreateProductParams{
			ProductRefNo: arg.ProductRefNo,
			ProductName: arg.ProductName,
			ProductDescription: arg.ProductDescription,
			ProductCode: arg.ProductCode,
			Price: arg.Price,
			SalePrice: arg.SalePrice,
			Collection: arg.Collection,
			Quantity: arg.Quantity,
			Color: arg.Color,
			Size: arg.Size,
			Status: arg.Status,
		})

		if err != nil {
			return err
		}

		// Main image for product
		mainImage, err := q.CreateProductMedia(ctx, CreateProductMediaParams{
			ProductID: result.Product.ProductRefNo,
			IsMainImage: true,
			ProductMediaRef: util.RandomString(12),
			MediaID: arg.MainImage,
		})

		if err != nil {
			return err
		}

		// Other Image 1 for product
		otherImage1, err := q.CreateProductMedia(ctx, CreateProductMediaParams{
			ProductID: result.Product.ProductRefNo,
			IsMainImage: false,
			ProductMediaRef: util.RandomString(12),
			MediaID: arg.OtherImage1,
		})

		if err != nil {
			return err
		}

		// Other Image 2 for product
		otherImage2, err := q.CreateProductMedia(ctx, CreateProductMediaParams{
			ProductID: result.Product.ProductRefNo,
			IsMainImage: false,
			ProductMediaRef: util.RandomString(12),
			MediaID: arg.OtherImage2,
		})

		if err != nil {
			return err
		}

		// Other Image 3 for product
		otherImage3, err := q.CreateProductMedia(ctx, CreateProductMediaParams{
			ProductID: result.Product.ProductRefNo,
			IsMainImage: false,
			ProductMediaRef: util.RandomString(12),
			MediaID: arg.OtherImage3,
		})

		if err != nil {
			return err
		}

		// update the product with the details
		result.Product, err = q.UpdateProduct(ctx, UpdateProductParams{
			ID: result.Product.ID,
			ProductImageMain: pgtype.Text{
				String: mainImage.ProductMediaRef,
				Valid: true,
			},
			ProductImageOther1: pgtype.Text{
				String: otherImage1.ProductMediaRef,
				Valid: true,
			},
			ProductImageOther2: pgtype.Text{
				String: otherImage2.ProductMediaRef,
				Valid: true,
			},
			ProductImageOther3: pgtype.Text{
				String: otherImage3.ProductMediaRef,
				Valid: true,
			},
		})

		if err != nil {
			return err
		}

		// finnally update the collection product count
		// first fetch the collection and retrieve the current quantity
		result.Collection, err = q.GetCollection(ctx, result.Product.Collection)

		if err != nil {
			return err
		}

		currentProductCount := result.Collection.ProductCount.Int64

		// Update the new product count
		result.Collection, err = q.UpdateCollection(ctx, UpdateCollectionParams{
			ID: result.Product.Collection,
			ProductCount: pgtype.Int8{
				Int64: currentProductCount + 1,
				Valid: true,
			},
		})

		return err
	})
	return result, err
}