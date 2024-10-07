package gapi

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateMediaAPI(t *testing.T) {
	media := randomMedia()

	testCases := []struct {
		name string
		req *pb.CreateMediaRequest
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.CreateMediaResponse, err error) 
	} {
		{
			name: "OK",
			req: &pb.CreateMediaRequest{
				MediaRef: media.MediaRef,
				Url: media.Url,
				AwsId: media.AwsID,
			},
			buildStubs: func (store *mockdb.MockStore)  {
				arg := db.CreateMediaParams {
					MediaRef: media.MediaRef,
					Url: media.Url,
					AwsID: media.AwsID,
				}

				store.EXPECT().
					CreateMedia(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(media, nil)
			},
			checkResponse:  func(t *testing.T, res *pb.CreateMediaResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				createdMedia := res.GetMedia()
				require.Equal(t, media.MediaRef, createdMedia.MediaRef)
				require.Equal(t, media.AwsID, createdMedia.AwsId)
				require.Equal(t, media.Url, createdMedia.Url)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateMediaRequest{
				MediaRef: media.MediaRef,
				Url: media.Url,
				AwsId: media.AwsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Simulate an internal error from the CreateMedia operation
				store.EXPECT().
					CreateMedia(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Medium{}, sql.ErrConnDone) // Return a database connection error
			},
			checkResponse: func(t *testing.T, res *pb.CreateMediaResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code()) // Check for the Internal error code
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

			server := newTestServer(t, store, nil)

			res, err := server.CreateMedia(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}