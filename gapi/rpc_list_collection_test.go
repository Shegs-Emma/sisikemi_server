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

func TestListCollectionAPI(t *testing.T) {
	user, _ := randomUser(t)
	collectionItems := []db.Collection{
		{
			ID: util.RandomInt(1, 1000), 
			CollectionName: util.RandomString(10), 
			CollectionDescription: util.RandomString(20),
			ProductCount: pgtype.Int8{
				Int64: 4,
				Valid: true,
			},
			ThumbnailImage: util.RandomString(15),
			HeaderImage: util.RandomString(15),
		},
		{
			ID: util.RandomInt(1, 1000), 
			CollectionName: util.RandomString(10), 
			CollectionDescription: util.RandomString(20),
			ProductCount: pgtype.Int8{
				Int64: 4,
				Valid: true,
			},
			ThumbnailImage: util.RandomString(15),
			HeaderImage: util.RandomString(15),
		},
		{
			ID: util.RandomInt(1, 1000), 
			CollectionName: util.RandomString(10), 
			CollectionDescription: util.RandomString(20),
			ProductCount: pgtype.Int8{
				Int64: 4,
				Valid: true,
			},
			ThumbnailImage: util.RandomString(15),
			HeaderImage: util.RandomString(15),
		},
	}

	testCases := []struct {
		name	string
		req		*pb.ListCollectionRequest
		buildStubs		func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse 	func(t *testing.T, res *pb.ListCollectionResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListCollectionRequest{
				PageSize: 3,
				PageId: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCollection(gomock.Any(), db.ListCollectionParams{
						Limit: 3,
						Offset: 0,
					}).
					Return(collectionItems, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListCollectionResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, len(collectionItems), len(res.Collection))
				for i, collection := range collectionItems {
					require.Equal(t, collection.ID, res.Collection[i].Id)
					require.Equal(t, collection.CollectionName, res.Collection[i].CollectionName)
					require.Equal(t, collection.CollectionDescription, res.Collection[i].CollectionDescription)
					require.Equal(t, collection.ProductCount.Int64, res.Collection[i].ProductCount)
					require.Equal(t, collection.ThumbnailImage, res.Collection[i].ThumbnailImage)
					require.Equal(t, collection.HeaderImage, res.Collection[i].HeaderImage)
				}
			},
		},
		{
			name: "NotFound",
			req: &pb.ListCollectionRequest{
				PageSize:  2,
				PageId:    1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCollection(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.ListCollectionResponse, err error) {
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
			res, err := server.ListCollections(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}