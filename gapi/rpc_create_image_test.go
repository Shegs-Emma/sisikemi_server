package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateImageAPI(t *testing.T) {
	testCases := []struct {
		name string
		req *pb.UploadImageRequest
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.UploadImageResponse, err error) 
	} {
		{
			name: "OK",
			req: &pb.UploadImageRequest{
				Image:    []byte(util.RandomString(20)),
				Filename: util.RandomString(12),
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Allow matching any generated values instead of strict equality
				store.EXPECT().
					CreateMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, arg db.CreateMediaParams) (db.Medium, error) {
						require.NotEmpty(t, arg.MediaRef)
						require.Contains(t, arg.Url, "http://localhost:8080/uploads/")
						require.NotEmpty(t, arg.AwsID)

						return db.Medium{
							MediaRef: arg.MediaRef,
							Url:      arg.Url,
							AwsID:    arg.AwsID,
						}, nil
					}).
					Times(1)
			},
			checkResponse: func(t *testing.T, res *pb.UploadImageResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				createdMedia := res.GetMedia()
				require.NotEmpty(t, createdMedia.MediaRef)
				require.Contains(t, createdMedia.Url, "http://localhost:8080/uploads/")
				require.NotEmpty(t, createdMedia.AwsId)
			},
		},
		{
			name: "InternalError",
			req: &pb.UploadImageRequest{
				Image:    []byte(util.RandomString(20)),
				Filename: util.RandomString(12),
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Simulate an internal error from the CreateMedia operation
				store.EXPECT().
					CreateMedia(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Medium{}, sql.ErrConnDone) // Return a database connection error
			},
			checkResponse: func(t *testing.T, res *pb.UploadImageResponse, err error) {
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

			server := newTestServer(t, store, nil, nil)

			res, err := server.UploadImage(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}

func randomImage() db.Medium {
	return db.Medium{
		ID: util.RandomInt(1, 1000),
		Url: fmt.Sprintf("http://localhost:8080/uploads/%s", util.RandomString(9)),
		AwsID: util.RandomString(15),
		MediaRef: util.RandomString(10),
	}
}