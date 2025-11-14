package gapi

import (
	"context"
	"database/sql"
	"testing"

	mockdb "github.com/Shegs-Emma/sisikemi_server/db/mock"
	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetMediaAPI(t *testing.T) {
	media := randomImage()

	testCases := []struct {
		name          string
		req           *pb.GetMediaRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.GetMediaResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetMediaRequest{
				Id: media.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMedia(gomock.Any(), gomock.Eq(media.ID)).
					Times(1).
					Return(media, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetMediaResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, media.ID, res.Media.Id)
				require.Equal(t, media.Url, res.Media.Url)
				require.Equal(t, media.AwsID, res.Media.AwsId)
				require.Equal(t, media.MediaRef, res.Media.MediaRef)
			},
		},
		{
			name: "NotFound",
			req: &pb.GetMediaRequest{
				Id: media.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMedia(gomock.Any(), gomock.Eq(media.ID)).
					Times(1).
					Return(db.Medium{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *pb.GetMediaResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.GetMediaRequest{
				Id: media.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMedia(gomock.Any(), gomock.Eq(media.ID)).
					Times(1).
					Return(db.Medium{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.GetMediaResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code()) // Change if your implementation changes
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

			res, err := server.GetMedia(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}
