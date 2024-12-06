package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type DeleteProductTxParams struct {
	ProductId int32
}

type DeleteProductTxResult struct {
	Product Product
	ProductMedium ProductMedium
	Collection Collection
}

func (store *SQLStore) DeleteProductTx (ctx context.Context, arg DeleteProductTxParams) (DeleteProductTxResult, error) {
	var result DeleteProductTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Fetch the product
		result.Product, err = q.GetProduct(ctx, int64(arg.ProductId))

		if err != nil {
			return err
		}

		// Delete the images associated with the product
		err = q.DeleteProductMedia(ctx, result.Product.ProductRefNo)

		if err != nil {
			return err
		}

		// Update product count in collection
		// fetch and retreive the current count
		result.Collection, err = q.GetCollection(ctx, result.Product.Collection)

		if err != nil {
			return err
		}

		currentProductCount := result.Collection.ProductCount.Int64

		// update the new product count
		result.Collection, err = q.UpdateCollection(ctx, UpdateCollectionParams{
			ID: result.Product.Collection,
			ProductCount: pgtype.Int8{
				Int64: currentProductCount - 1,
				Valid: true,
			},
		})

		if err != nil {
			return err
		}

		// Delete the product
		err = q.DeleteProduct(ctx, int64(arg.ProductId))

		if err != nil {
			return err
		}

		return err
	})

	return result, err
}