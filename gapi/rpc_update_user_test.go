package gapi

import (
	"context"
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

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	media := randomMedia()

	newFirstName := util.RandomUser()
	newLastName := util.RandomUser()
	newEmail := util.RandomEmail()
	newPhoneNumber := util.RandomString(11)
	newProfilePhoto := media.MediaRef
	invalidEmail := "invalid-email"

	testCases := []struct {
		name          string
		req          *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username:  user.Username,
				FirstName: &newFirstName,
				LastName: &newLastName,
				PhoneNumber: &newPhoneNumber,
				ProfilePhoto: &newProfilePhoto,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FirstName: pgtype.Text{
						String: newFirstName,
						Valid: true,
					},
					LastName: pgtype.Text{
						String: newLastName,
						Valid: true,
					},
					PhoneNumber: pgtype.Text{
						String: newPhoneNumber,
						Valid: true,
					},
					ProfilePhoto: pgtype.Text{
						String: newProfilePhoto,
						Valid: true,
					},
				}

				updatedUser := db.User{
					Username: user.Username,
					HashedPassword: user.HashedPassword,
					FirstName: newFirstName,
					LastName: newLastName,
					PhoneNumber: newPhoneNumber,
					ProfilePhoto: arg.ProfilePhoto,
					IsAdmin: user.IsAdmin,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt: user.CreatedAt,
					IsEmailVerified: user.IsEmailVerified,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()

				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newFirstName, updatedUser.FirstName)
				require.Equal(t, newLastName, updatedUser.LastName)
				require.Equal(t, newPhoneNumber, updatedUser.PhoneNumber)
				require.Equal(t, newProfilePhoto, updatedUser.ProfilePhoto)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username:  user.Username,
				FirstName: &newFirstName,
				LastName: &newLastName,
				PhoneNumber: &newPhoneNumber,
				ProfilePhoto: &newProfilePhoto,
				Email:     &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username:  user.Username,
				FirstName: &newFirstName,
				LastName: &newLastName,
				PhoneNumber: &newPhoneNumber,
				ProfilePhoto: &newProfilePhoto,
				Email:     &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.UpdateUserRequest{
				Username:  user.Username,
				FirstName: &newFirstName,
				LastName: &newLastName,
				PhoneNumber: &newPhoneNumber,
				ProfilePhoto: &newProfilePhoto,
				Email:     &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Username:  user.Username,
				FirstName: &newFirstName,
				LastName: &newLastName,
				PhoneNumber: &newPhoneNumber,
				ProfilePhoto: &newProfilePhoto,
				Email:     &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.IsAdmin, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
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
			res, err := server.UpdateUser(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}