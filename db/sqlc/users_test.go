package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T, media Medium) User {
	arg := CreateUserParams{
		Username: util.RandomUser(),
		HashedPassword: "secret",
		FirstName: util.RandomUser(),
		LastName: util.RandomUser(),
		Email: util.RandomEmail(),
		PhoneNumber: util.RandomString(11),
		ProfilePhoto: media.MediaRef,
		IsAdmin: false,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.ProfilePhoto, user.ProfilePhoto)
	require.Equal(t, arg.IsAdmin, user.IsAdmin)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	media := createRandomMedia(t)
	createRandomUser(t, media)
}

func TestGetUser(t *testing.T) {
	media := createRandomMedia(t)

	user1 := createRandomUser(t, media)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.ProfilePhoto, user2.ProfilePhoto)
	require.Equal(t, user1.IsAdmin, user2.IsAdmin)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

// func TestGetUserForUpdate(t *testing.T) {
// 	media := createRandomMedia(t)

// 	user1 := createRandomUser(t, media)
// 	user2, err := testQueries.GetUserForUpdate(context.Background(), user1.ID)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, user2)
// 	require.Equal(t, user1.Username, user2.Username)
// 	require.Equal(t, user1.FirstName, user2.FirstName)
// 	require.Equal(t, user1.LastName, user2.LastName)
// 	require.Equal(t, user1.Email, user2.Email)
// 	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
// 	require.Equal(t, user1.ProfilePhoto, user2.ProfilePhoto)
// 	require.Equal(t, user1.IsAdmin, user2.IsAdmin)
// 	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
// 	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
// 	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
// }

// func TestUpdateUser(t *testing.T) {
// 	media := createRandomMedia(t)

// 	user1 := createRandomUser(t, media)

// 	arg := UpdateUsersParams{
// 		ID: user1.ID,
// 		Username: util.RandomUser(),
// 		HashedPassword: "secret",
// 		FirstName: util.RandomUser(),
// 		LastName: util.RandomUser(),
// 		Email: util.RandomEmail(),
// 		PhoneNumber: util.RandomString(11),
// 		ProfilePhoto: media.MediaRef,
// 		IsAdmin: true,
// 	}

// 	user2, err := testQueries.UpdateUsers(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, user2)
// 	require.Equal(t, user1.ID, user2.ID)
// 	require.Equal(t, arg.Username, user2.Username)
// 	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
// 	require.Equal(t, arg.FirstName, user2.FirstName)
// 	require.Equal(t, arg.LastName, user2.LastName)
// 	require.Equal(t, arg.Email, user2.Email)
// 	require.Equal(t, arg.PhoneNumber, user2.PhoneNumber)
// 	require.Equal(t, arg.ProfilePhoto, user2.ProfilePhoto)
// 	require.Equal(t, arg.IsAdmin, user2.IsAdmin)
// 	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
// }

// func TestDeleteUser(t *testing.T){
// 	media := createRandomMedia(t)

// 	user1 := createRandomUser(t, media)
// 	err := testQueries.DeleteUser(context.Background(), user1.ID)

// 	require.NoError(t, err)
// 	user2, err := testQueries.GetUser(context.Background(), user1.Username)

// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, user2)
// }

// func TestListUsers(t *testing.T){
// 	for i := 0; i < 5; i++ {
// 		media := createRandomMedia(t)
// 		createRandomUser(t, media)
// 	}

// 	arg := ListUsersParams{
// 		Limit: 3,
// 		Offset: 3,
// 	}
// 	users, err := testQueries.ListUsers(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, users)
// 	require.Len(t, users, 3)

// 	for _, user := range users {
// 		require.NotEmpty(t, user)
// 	}
// }