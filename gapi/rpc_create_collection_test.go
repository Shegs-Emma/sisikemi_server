package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateCollectionAPI(t *testing.T) {
	user, _ := randomUser(t)
	collection := randomCollection()

	testCases := []struct {
		name string
		req *pb.CreateCollectionRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateCollectionResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateCollectionRequest{
				CollectionName: collection.CollectionName,
				CollectionDescription: collection.CollectionDescription,
				ThumbnailImage: collection.ThumbnailImage,
				HeaderImage: collection.HeaderImage,
			},
			buildStubs: func (store *mockdb.MockStore) {
				arg := db.CreateCollectionParams{
					CollectionName: collection.CollectionName,
					CollectionDescription: collection.CollectionDescription,
					ThumbnailImage: collection.ThumbnailImage,
					HeaderImage: collection.HeaderImage,
				}	

				createdCollection := db.Collection{
					CollectionName: collection.CollectionName,
					CollectionDescription: collection.CollectionDescription,
					ProductCount: pgtype.Int8{
						Int64: 7, // This is determined by the backend
						Valid: true,
					},
					ThumbnailImage: collection.ThumbnailImage,
					HeaderImage: collection.HeaderImage,
				}

				store.EXPECT().
					CreateCollection(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(createdCollection, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCollectionResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdCollection := res.GetCollection()

				require.Equal(t, collection.CollectionName, createdCollection.CollectionName)
				require.Equal(t, collection.CollectionDescription, createdCollection.CollectionDescription)
				require.Equal(t, collection.ProductCount.Int64, createdCollection.ProductCount)
				require.Equal(t, collection.ThumbnailImage, createdCollection.ThumbnailImage)
				require.Equal(t, collection.HeaderImage, createdCollection.HeaderImage)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateCollectionRequest{
				CollectionName: collection.CollectionName,
				CollectionDescription: collection.CollectionDescription,
				ThumbnailImage: collection.ThumbnailImage,
				HeaderImage: collection.HeaderImage,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Simulate an internal error from the CreateMedia operation
				store.EXPECT().
					CreateCollection(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Collection{}, sql.ErrConnDone) // Return a database connection error
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCollectionResponse, err error) {
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
			req: &pb.CreateCollectionRequest{
				CollectionName: collection.CollectionName,
				CollectionDescription: collection.CollectionDescription,
				ThumbnailImage: collection.ThumbnailImage,
				HeaderImage: collection.HeaderImage,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					CreateCollection(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.CreateCollectionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NotAdmin",
			req: &pb.CreateCollectionRequest{
				CollectionName: collection.CollectionName,
				CollectionDescription: collection.CollectionDescription,
				ThumbnailImage: collection.ThumbnailImage,
				HeaderImage: collection.HeaderImage,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
				CreateCollection(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, false, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateCollectionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.PermissionDenied, st.Code())
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
			res, err := server.CreateCollection(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}

func randomCollection() db.Collection{
	return db.Collection{
		ID: util.RandomInt(1, 1000),
		CollectionName: util.RandomString(7),
		CollectionDescription: util.RandomString(20),
		ProductCount: pgtype.Int8{
            Int64:  7, 
            Valid:  true, // `Valid` is used instead of `Status`
        },
		ThumbnailImage: util.RandomString(12),
		HeaderImage: util.RandomString(12),

	}
}