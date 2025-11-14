package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/Shegs-Emma/sisikemi_server/db/mock"
	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/token"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestListProductAPI(t *testing.T) {
	user, _ := randomUser(t)
	collection := randomCollection()
	productItems := []db.Product{
		{
			ID: util.RandomInt(1, 1000), 
			ProductRefNo: util.RandomString(10), 
			ProductName: util.RandomString(10),
			ProductDescription: util.RandomString(20),
			ProductCode: util.RandomString(3),
			Price: util.RandomInt(1, 20),
			SalePrice: util.RandomString(5),
			ProductImageMain: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther1: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther2: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther3:pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			Collection: collection.ID,
			Quantity: int32(util.RandomInt(1, 20)),
			Color: []string{"Yellow", "Red"},
			Size: []string{"XXL"},
			Status: db.ProductStatusActive,
		},
		{
			ID: util.RandomInt(1, 1000), 
			ProductRefNo: util.RandomString(10), 
			ProductName: util.RandomString(10),
			ProductDescription: util.RandomString(20),
			ProductCode: util.RandomString(3),
			Price: util.RandomInt(1, 20),
			SalePrice: util.RandomString(5),
			ProductImageMain: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther1: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther2: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther3:pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			Collection: collection.ID,
			Quantity: int32(util.RandomInt(1, 20)),
			Color: []string{"Yellow", "Red"},
			Size: []string{"XXL"},
			Status: db.ProductStatusActive,
		},
		{
			ID: util.RandomInt(1, 1000), 
			ProductRefNo: util.RandomString(10), 
			ProductName: util.RandomString(10),
			ProductDescription: util.RandomString(20),
			ProductCode: util.RandomString(3),
			Price: util.RandomInt(1, 20),
			SalePrice: util.RandomString(5),
			ProductImageMain: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther1: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther2: pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			ProductImageOther3:pgtype.Text{
				String: util.RandomString(10),
				Valid: true,
			},
			Collection: collection.ID,
			Quantity: int32(util.RandomInt(1, 20)),
			Color: []string{"Yellow", "Red"},
			Size: []string{"XXL"},
			Status: db.ProductStatusActive,
		},
	}

	testCases := []struct {
		name	string
		req		*pb.ListProductRequest
		buildStubs		func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse 	func(t *testing.T, res *pb.ListProductResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListProductRequest{
				PageSize: 3,
				PageId: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Any()).
					Return(db.ProductMedium{}, sql.ErrNoRows).
					AnyTimes()

				store.EXPECT().
					GetCollection(gomock.Any(), gomock.Any()).
					Return(db.Collection{}, sql.ErrNoRows).
					AnyTimes()

				store.EXPECT().
					ListProducts(gomock.Any(), db.ListProductsParams{
						Limit: 3,
						Offset: 0,
					}).
					Return(productItems, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, len(productItems), len(res.Product))
				for i, product := range productItems {
					require.Equal(t, product.ID, res.Product[i].Id)
					require.Equal(t, product.ProductName, res.Product[i].ProductName)
					require.Equal(t, product.ProductDescription, res.Product[i].ProductDescription)
					require.Equal(t, product.ProductRefNo, res.Product[i].ProductRefNo)
					require.Equal(t, product.ProductCode, res.Product[i].ProductCode)
					require.Equal(t, product.Color, res.Product[i].Color)
					require.Equal(t, product.Size, res.Product[i].Size)
					require.Equal(t, product.Price, res.Product[i].Price)
				}
			},
		},
		{
			name: "NotFound",
			req: &pb.ListProductRequest{
				PageSize:  2,
				PageId:    1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListProductResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil)
			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.ListProducts(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}