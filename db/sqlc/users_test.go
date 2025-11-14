package db

import (
	"context"
	"testing"
	"time"

	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, media Medium) User {
	arg := CreateUserParams{
		Username: util.RandomUser(),
		HashedPassword: "secret",
		FirstName: util.RandomUser(),
		LastName: util.RandomUser(),
		Email: util.RandomEmail(),
		PhoneNumber: util.RandomString(11),
		ProfilePhoto: pgtype.Text{
			String: media.MediaRef,
			Valid: true,
		},
		IsAdmin: false,
	}

	user, err := testStore.CreateUser(context.Background(), arg)

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
	user2, err := testStore.GetUserByUsername(context.Background(), user1.Username)

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

func TestUpdateUserOnlyFirstName(t *testing.T) {
	media := createRandomMedia(t)
	oldUser := createRandomUser(t, media)

	newFirstName := util.RandomUser()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FirstName: pgtype.Text{
			String: newFirstName,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
}

func TestUpdateUserOnlyLastName(t *testing.T) {
	media := createRandomMedia(t)
	oldUser := createRandomUser(t, media)

	newLastName := util.RandomUser()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		LastName: pgtype.Text{
			String: newLastName,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
}

func TestUpdateUserOnlyPhoneNumber(t *testing.T) {
	media := createRandomMedia(t)
	oldUser := createRandomUser(t, media)

	newPhoneNumber := util.RandomUser()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		PhoneNumber: pgtype.Text{
			String: newPhoneNumber,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
}

func TestUpdateUserOnlyProfilePhoto(t *testing.T) {
	oldMedia := createRandomMedia(t)
	oldUser := createRandomUser(t, oldMedia)

	newMedia := createRandomMedia(t)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		ProfilePhoto: pgtype.Text{
			String: newMedia.MediaRef,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	media := createRandomMedia(t)
	oldUser := createRandomUser(t, media)

	newPassword := util.RandomString(8)
	newHashPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashPassword,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
}

func TestUpdateUserAllFields(t *testing.T) {
	media := createRandomMedia(t)
	oldUser := createRandomUser(t, media)

	newFirstName := util.RandomUser()
	newLastName := util.RandomUser()
	newPhoneNumber := util.RandomString(11)
	newMedia := createRandomMedia(t)
	newPassword := util.RandomString(8)
	newHashPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashPassword,
			Valid: true,
		},
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
			String: newMedia.MediaRef,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashPassword)

	require.NotEqual(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, updatedUser.FirstName, newFirstName)

	require.NotEqual(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, updatedUser.LastName, newLastName)

	require.NotEqual(t, oldUser.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, updatedUser.PhoneNumber, newPhoneNumber)

	require.NotEqual(t, oldUser.ProfilePhoto, updatedUser.ProfilePhoto)
	require.Equal(t, updatedUser.ProfilePhoto.String, newMedia.MediaRef)
}