package gapi

import (
	"context"
	"database/sql"
	"testing"

	mockdb "github.com/Shegs-Emma/sisikemi_server/db/mock"
	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListMediaAPI(t *testing.T) {
	mediaItems := []db.Medium{
		{ID: util.RandomInt(1, 1000), MediaRef: util.RandomString(10), Url: util.RandomString(10), AwsID: util.RandomString(15)},
		{ID: util.RandomInt(1, 1000), MediaRef: util.RandomString(10), Url: util.RandomString(10), AwsID: util.RandomString(15)},
	}

	testCases := []struct {
		name          string
		req           *pb.ListMediaRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.ListMediaResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ListMediaRequest{
				PageSize: 2,
				PageId: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMedia(gomock.Any(), db.ListMediaParams{
						Limit:  2,
						Offset: 0,
					}).
					Return(mediaItems, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListMediaResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, len(mediaItems), len(res.Media))
				for i, media := range mediaItems {
					require.Equal(t, media.ID, res.Media[i].Id)
					require.Equal(t, media.MediaRef, res.Media[i].MediaRef)
					require.Equal(t, media.Url, res.Media[i].Url)
					require.Equal(t, media.AwsID, res.Media[i].AwsId)
				}
			},
		},
		{
			name: "NotFound",
			req: &pb.ListMediaRequest{
				PageSize:  2,
				PageId:    1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMedia(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *pb.ListMediaResponse, err error) {
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

			res, err := server.ListMedia(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}
