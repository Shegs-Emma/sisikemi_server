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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateCartItemAPI(t *testing.T) {
	user, _ := randomUser(t)
	cartItem := randomCartItem(user)

	testCases := []struct {
		name string
		req *pb.CreateCartItemRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateCartItemResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateCartItemRequest{
				ProductId: int64(cartItem.ProductID),
				ProductName: cartItem.ProductName,
				ProductPrice: cartItem.ProductPrice,
				ProductQuantity: cartItem.ProductQuantity,
				ProductImage: cartItem.ProductImage,
				ProductColor: cartItem.ProductColor,
				ProductSize: cartItem.ProductSize,
			},
			buildStubs: func (store *mockdb.MockStore) {
				arg := db.CreateCartItemParams{
					ProductID: cartItem.ProductID,
					ProductName: cartItem.ProductName,
					ProductPrice: cartItem.ProductPrice,
					ProductQuantity: cartItem.ProductQuantity,
					ProductImage: cartItem.ProductImage,
					ProductColor: cartItem.ProductColor,
					ProductSize: cartItem.ProductSize,
					UserRefID: cartItem.UserRefID,
				}

				createdCartItem := db.Cart{
					ProductID: cartItem.ProductID,
					ProductName: cartItem.ProductName,
					ProductPrice: cartItem.ProductPrice,
					ProductQuantity: cartItem.ProductQuantity,
					ProductImage: cartItem.ProductImage,
					ProductColor: cartItem.ProductColor,
					ProductSize: cartItem.ProductSize,
					UserRefID: cartItem.UserRefID,
				}

				store.EXPECT().
					GetCartItemByProductId(gomock.Any(), gomock.Eq(cartItem.ProductID)).
					Return(db.Cart{}, sql.ErrNoRows) 

				// Return the user object when fetching by ID
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(cartItem.UserRefID)).
					Return(user, nil)

				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Return(user, nil)

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(createdCartItem, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCartItemResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdCart := res.GetCart()

				require.Equal(t, cartItem.ProductName, createdCart.ProductName)
				require.Equal(t, cartItem.ProductPrice, createdCart.ProductPrice)
				require.Equal(t, cartItem.ProductQuantity, createdCart.ProductQuantity)
				require.Equal(t, cartItem.ProductImage, createdCart.ProductImage)
				require.Equal(t, cartItem.ProductColor, createdCart.ProductColor)
				require.Equal(t, cartItem.ProductSize, createdCart.ProductSize)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateCartItemRequest{
				ProductId: int64(cartItem.ProductID),
				ProductName: cartItem.ProductName,
				ProductPrice: cartItem.ProductPrice,
				ProductQuantity: cartItem.ProductQuantity,
				ProductImage: cartItem.ProductImage,
				ProductColor: cartItem.ProductColor,
				ProductSize: cartItem.ProductSize,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCartItemByProductId(gomock.Any(), gomock.Eq(cartItem.ProductID)).
					Return(db.Cart{}, sql.ErrNoRows) 
				
				// store.EXPECT().
				// 	GetUser(gomock.Any(), gomock.Eq(user.Username)).
				// 	Return(user, nil) 

				// Return the user object when fetching by ID
				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(int64(user.ID))).
					Return(db.User{}, sql.ErrNoRows).AnyTimes()

				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Return(db.User{}, nil)

				// Simulate an internal error from the CreateMedia operation
				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Cart{}, sql.ErrConnDone) // Return a database connection error
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCartItemResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code()) // Check for the Internal error code
				// Ensure that the response is nil when an error occurs
				require.Nil(t, res)
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.CreateCartItemRequest{
				ProductId: int64(cartItem.ProductID),
				ProductName: cartItem.ProductName,
				ProductPrice: cartItem.ProductPrice,
				ProductQuantity: cartItem.ProductQuantity,
				ProductImage: cartItem.ProductImage,
				ProductColor: cartItem.ProductColor,
				ProductSize: cartItem.ProductSize,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateCartItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.CreateCartItemResponse, err error) {
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
			res, err := server.CreateCartItem(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}

func randomCartItem(user db.User) db.Cart{
	return db.Cart{
		ID: util.RandomInt(1, 1000),
		ProductID: int32(util.RandomInt(1, 1000)),
		ProductName: util.RandomString(20),
		ProductPrice: util.RandomInt(1, 1000),
		ProductQuantity: util.RandomInt(1, 1000),
		ProductImage: util.RandomString(20),
		ProductColor: "Yellow",
		ProductSize: "XXL",
		UserRefID: user.ID,
	}
}