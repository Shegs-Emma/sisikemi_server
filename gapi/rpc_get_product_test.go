package gapi

import (
	"context"
	"testing"

	mockdb "github.com/Shegs-Emma/sisikemi_server/db/mock"
	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestGetProductAPI(t *testing.T) {
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
		name          string
		req           *pb.GetProductRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.GetProductResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetProductRequest{
				Id: product.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCollection(gomock.Any(), gomock.Eq(product.Collection)).
					Return(collection, nil) 
				
				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Eq(product.ProductImageMain.String)).
					Return(media1, nil)
				

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Eq(product.ProductImageOther1.String)).
					Return(media2, nil)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Eq(product.ProductImageOther2.String)).
					Return(media3, nil)

				store.EXPECT().
					GetProductMediaByRef(gomock.Any(), gomock.Eq(product.ProductImageOther3.String)).
					Return(media4, nil)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Eq(mediaHeader.MediaRef)).
					Return(mediaHeader, nil)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Eq(mediaOther1.MediaRef)).
					Return(mediaOther1, nil)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Eq(mediaOther2.MediaRef)).
					Return(mediaOther2, nil)

				store.EXPECT().
					GetMediaByRef(gomock.Any(), gomock.Eq(mediaOther3.MediaRef)).
					Return(mediaOther3, nil)
				
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetProductResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, product.ID, res.Product.Id)
				require.Equal(t, product.ProductRefNo, res.Product.ProductRefNo)
				require.Equal(t, product.ProductName, res.Product.ProductName)
				require.Equal(t, product.ProductDescription, res.Product.ProductDescription)
				require.Equal(t, product.ProductCode, res.Product.ProductCode)
				require.Equal(t, product.Price, res.Product.Price)
				require.Equal(t, product.SalePrice, res.Product.SalePrice)

				require.Equal(t, product.Quantity, res.Product.Quantity)
				require.Equal(t, product.Color, res.Product.Color)
				require.Equal(t, product.Size, res.Product.Size)
				require.Equal(t, string(product.Status), res.Product.Status)
				require.Equal(t, product.ProductImageMain.String, res.Product.ProductImageMain.ProductMediaRef)
				require.Equal(t, product.ProductImageMain.Valid, true)

				require.Equal(t, product.ProductImageOther1.String, res.Product.ProductImageOther_1.ProductMediaRef)
				require.Equal(t, product.ProductImageOther1.Valid, true)

				require.Equal(t, product.ProductImageOther2.String, res.Product.ProductImageOther_2.ProductMediaRef)
				require.Equal(t, product.ProductImageOther2.Valid, true)

				require.Equal(t, product.ProductImageOther3.String, res.Product.ProductImageOther_3.ProductMediaRef)
				require.Equal(t, product.ProductImageOther3.Valid, true)

			},
		},
		// {
		// 	name: "NotFound",
		// 	req: &pb.GetMediaRequest{
		// 		Id: media.ID,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetMedia(gomock.Any(), gomock.Eq(media.ID)).
		// 			Times(1).
		// 			Return(db.Medium{}, sql.ErrNoRows)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.GetMediaResponse, err error) {
		// 		require.Error(t, err)
		// 		st, ok := status.FromError(err)
		// 		require.True(t, ok)
		// 		require.Equal(t, codes.NotFound, st.Code())
		// 	},
		// },
		// {
		// 	name: "InternalError",
		// 	req: &pb.GetMediaRequest{
		// 		Id: media.ID,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetMedia(gomock.Any(), gomock.Eq(media.ID)).
		// 			Times(1).
		// 			Return(db.Medium{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, res *pb.GetMediaResponse, err error) {
		// 		require.Error(t, err)
		// 		st, ok := status.FromError(err)
		// 		require.True(t, ok)
		// 		require.Equal(t, codes.NotFound, st.Code()) // Change if your implementation changes
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil)

			res, err := server.GetProduct(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}

func randomProduct(collection db.Collection, media1 db.ProductMedium, media2 db.ProductMedium, media3 db.ProductMedium, media4 db.ProductMedium) db.Product {
	return db.Product{
		ID: util.RandomInt(1, 1000),
		ProductRefNo: util.RandomString(12),
		ProductName: util.RandomString(15),
		ProductDescription: util.RandomString(20),
		ProductCode: util.RandomString(4),
		Price: util.RandomInt(1, 1000),
		SalePrice: util.RandomString(4),
		Collection: collection.ID,
		Quantity: int32(util.RandomInt(1, 1000)),
		Color: []string{"Yellow", "Red"},
		Size: []string{"S", "XXL"},
		Status: "active",
		ProductImageMain: pgtype.Text{
            String:  media1.ProductMediaRef, 
            Valid:  true, // `Valid` is used instead of `Status`
        },
		ProductImageOther1: pgtype.Text{
            String:  media2.ProductMediaRef, 
            Valid:  true, // `Valid` is used instead of `Status`
        },
		ProductImageOther2: pgtype.Text{
            String:  media3.ProductMediaRef, 
            Valid:  true, // `Valid` is used instead of `Status`
        },
		ProductImageOther3: pgtype.Text{
            String:  media4.ProductMediaRef, 
            Valid:  true, // `Valid` is used instead of `Status`
        },
	}
}

func randomProductMedia(media db.Medium) db.ProductMedium {
	return db.ProductMedium{
		ID: util.RandomInt(1, 1000),
		ProductMediaRef: util.RandomString(10),
		ProductID: util.RandomString(15),
		MediaID: media.MediaRef,
		IsMainImage: false,
	}
}