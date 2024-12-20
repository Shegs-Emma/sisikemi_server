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
	"github.com/techschool/simplebank/util"
)

func TestListUserCartItemsAPI(t *testing.T) {
	user, _ := randomUser(t)
	userCartItems := []db.Cart{
		{
			ID: util.RandomInt(1, 1000), 
			ProductID: int32(util.RandomInt(1, 10)), 
			ProductName: util.RandomString(20),
			ProductPrice: util.RandomString(20),
			ProductQuantity: util.RandomInt(1, 10),
			ProductImage: util.RandomString(10),
			ProductColor: "Yellow",
			ProductSize: "XL",
			UserRefID: int64(user.ID),
		},
		{
			ID: util.RandomInt(1, 1000), 
			ProductID: int32(util.RandomInt(1, 10)), 
			ProductName: util.RandomString(20),
			ProductPrice: util.RandomString(20),
			ProductQuantity: util.RandomInt(1, 10),
			ProductImage: util.RandomString(10),
			ProductColor: "Yellow",
			ProductSize: "XL",
			UserRefID: int64(user.ID),
		},
		{
			ID: util.RandomInt(1, 1000), 
			ProductID: int32(util.RandomInt(1, 10)), 
			ProductName: util.RandomString(20),
			ProductPrice: util.RandomString(20),
			ProductQuantity: util.RandomInt(1, 10),
			ProductImage: util.RandomString(10),
			ProductColor: "Yellow",
			ProductSize: "XL",
			UserRefID: int64(user.ID),
		},
	}

	testCases := []struct {
		name	string
		req		*pb.ListUserCartRequest
		buildStubs		func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse 	func(t *testing.T, res *pb.ListUserCartResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListUserCartRequest{
				PageSize: 3,
				PageId: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Return(user, nil).AnyTimes()


				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(user, nil).AnyTimes()
					

				store.EXPECT().
					ListUserCartItems(gomock.Any(), db.ListUserCartItemsParams{
						UserRefID: user.ID,
						Limit: 3,
						Offset: 0,
					}).
					Return(userCartItems, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListUserCartResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, len(userCartItems), len(res.Cart))
				for i, userCart := range userCartItems {
					require.Equal(t, userCart.ID, res.Cart[i].Id)
					require.Equal(t, int64(userCart.UserRefID), int64(user.ID))
					require.Equal(t, userCart.ProductColor, res.Cart[i].ProductColor)
					require.Equal(t, userCart.ProductName, res.Cart[i].ProductName)
					require.Equal(t, userCart.ProductID, int32(res.Cart[i].ProductId))
					require.Equal(t, userCart.ProductPrice, res.Cart[i].ProductPrice)
					require.Equal(t, userCart.ProductQuantity, res.Cart[i].ProductQuantity)
					require.Equal(t, userCart.ProductSize, res.Cart[i].ProductSize)
				}
			},
		},
		{
			name: "NotFound",
			req: &pb.ListUserCartRequest{
				PageSize:  2,
				PageId:    1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Return(user, nil).AnyTimes()


				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Any()).
					Return(user, nil).AnyTimes()
					
				store.EXPECT().
					ListUserCartItems(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListUserCartResponse, err error) {
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
			res, err := server.ListUserCartItems(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}