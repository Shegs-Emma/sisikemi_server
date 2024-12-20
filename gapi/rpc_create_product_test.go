package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateProductAPI(t *testing.T) {
	user, _ := randomUser(t)
	collection := randomCollection()
	mediaHeader := randomMedia()
	mediaOther1 := randomMedia()
	mediaOther2 := randomMedia()
	mediaOther3 := randomMedia()
	media1 := randomProductMedia(mediaHeader)
	media2 := randomProductMedia(mediaOther1)
	media3 := randomProductMedia(mediaOther2)
	media4 := randomProductMedia(mediaOther3)
	product := randomProduct(collection, media1, media2, media3, media4)

	testCases := []struct {
		name string
		req *pb.CreateProductRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateProductResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateProductRequest{
				ProductName: product.ProductName,
				ProductDescription: product.ProductDescription,
				ProductCode: product.ProductCode,
				Price: product.Price,
				SalePrice: product.SalePrice,
				Size: product.Size,
				Color: product.Color,
				Collection: product.Collection,
				Quantity: product.Quantity,
				Status: string(product.Status),
				MainImage: product.ProductImageMain.String,
				OtherImage_1: product.ProductImageOther1.String,
				OtherImage_2: product.ProductImageOther2.String,
				OtherImage_3: product.ProductImageOther3.String,
			},
			buildStubs: func (store *mockdb.MockStore) {
				store.EXPECT().
					GetCollection(gomock.Any(), gomock.Any()).
					Return(collection, nil).
					Times(1)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Any()).
					Return(mediaHeader, nil).
					Times(1)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Any()).
					Return(mediaOther1, nil).
					Times(1)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Any()).
					Return(mediaOther2, nil).
					Times(1)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Any()).
					Return(mediaOther3, nil).
					Times(1)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Any()).
					Return(media1, nil).
					Times(1)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Any()).
					Return(media2, nil).
					Times(1)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Any()).
					Return(media3, nil).
					Times(1)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Any()).
					Return(media4, nil).
					Times(1)

				store.EXPECT().
					CreateProductTx(gomock.Any(), gomock.AssignableToTypeOf(db.CreateProductTxParams{})).
					Times(1).
					Return(db.CreateProductTxResult{
						Product: product,
					}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdProduct := res.GetProduct()

				require.Equal(t, product.ProductRefNo, createdProduct.ProductRefNo)
				require.Equal(t, product.ProductName, createdProduct.ProductName)
				require.Equal(t, product.ProductDescription, createdProduct.ProductDescription)
				require.Equal(t, product.ProductCode, createdProduct.ProductCode)
				require.Equal(t, product.Price, createdProduct.Price)
				require.Equal(t, product.SalePrice, createdProduct.SalePrice)
				require.Equal(t, product.Size, createdProduct.Size)
				require.Equal(t, product.Color, createdProduct.Color)
				// require.Equal(t, product.Collection, createdProduct.Collection.Id)
				require.Equal(t, product.Quantity, createdProduct.Quantity)

				// require.Equal(t, product.Status, createdProduct.Status)
				// require.Equal(t, product.ProductImageMain, createdProduct.ProductImageMain)
				// require.Equal(t, product.ProductImageOther1, createdProduct.ProductImageOther_1)
				// require.Equal(t, product.ProductImageOther2, createdProduct.ProductImageOther_2)
				// require.Equal(t, product.ProductImageOther3, createdProduct.ProductImageOther_3)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateProductRequest{
				ProductName: product.ProductName,
				ProductDescription: product.ProductDescription,
				ProductCode: product.ProductCode,
				Price: product.Price,
				SalePrice: product.SalePrice,
				Size: product.Size,
				Color: product.Color,
				Collection: product.Collection,
				Quantity: product.Quantity,
				Status: string(product.Status),
				MainImage: product.ProductImageMain.String,
				OtherImage_1: product.ProductImageOther1.String,
				OtherImage_2: product.ProductImageOther2.String,
				OtherImage_3: product.ProductImageOther3.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateProductTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateProductTxResult{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateProductResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
				require.Nil(t, res)
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.CreateProductRequest{
				ProductName: product.ProductName,
				ProductDescription: product.ProductDescription,
				ProductCode: product.ProductCode,
				Price: product.Price,
				SalePrice: product.SalePrice,
				Size: product.Size,
				Color: product.Color,
				Collection: product.Collection,
				Quantity: product.Quantity,
				Status: string(product.Status),
				MainImage: product.ProductImageMain.String,
				OtherImage_1: product.ProductImageOther1.String,
				OtherImage_2: product.ProductImageOther2.String,
				OtherImage_3: product.ProductImageOther3.String,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateProductTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.CreateProductResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
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
			res, err := server.CreateProduct(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}