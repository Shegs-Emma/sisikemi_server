package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomMedia(t *testing.T) Medium {
	arg := CreateMediaParams {
		MediaRef: util.RandomString(6),
		Url: util.RandomString(15),
		AwsID: util.RandomString(10),
	}

	media, err := testStore.CreateMedia(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, media)
	require.Equal(t, arg.MediaRef, media.MediaRef)
	require.Equal(t, arg.Url, media.Url)
	require.Equal(t, arg.AwsID, media.AwsID)
	require.NotZero(t, media.ID)
	require.NotZero(t, media.CreatedAt)

	return media
}

func TestCreateMedia(t *testing.T) {
	createRandomMedia(t)
}

func TestGetMedia(t *testing.T) {
	media1 := createRandomMedia(t)
	media2, err := testStore.GetMedia(context.Background(), media1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, media2)
	require.Equal(t, media1.ID, media2.ID)
	require.Equal(t, media1.AwsID, media2.AwsID)
	require.Equal(t, media1.MediaRef, media2.MediaRef)
	require.Equal(t, media1.Url, media2.Url)
	require.WithinDuration(t, media1.CreatedAt, media2.CreatedAt, time.Second)
}

func TestGetMediaForUpdate(t *testing.T) {
	media1 := createRandomMedia(t)
	media2, err := testStore.GetMediaForUpdate(context.Background(), media1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, media2)
	require.Equal(t, media1.ID, media2.ID)
	require.Equal(t, media1.AwsID, media2.AwsID)
	require.Equal(t, media1.MediaRef, media2.MediaRef)
	require.Equal(t, media1.Url, media2.Url)
	require.WithinDuration(t, media1.CreatedAt, media2.CreatedAt, time.Second)
}

func TestUpdateMedia(t *testing.T) {
	media1 := createRandomMedia(t)
	arg := UpdateMediaParams{
		ID: media1.ID,
		Url: util.RandomString(15),
		AwsID: util.RandomString(10),
	}

	media2, err := testStore.UpdateMedia(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, media2)
	require.Equal(t, media1.ID, media2.ID)
	require.Equal(t, media1.MediaRef, media2.MediaRef)
	require.Equal(t, arg.Url, media2.Url)
	require.Equal(t, arg.AwsID, media2.AwsID)
	require.WithinDuration(t, media1.CreatedAt, media2.CreatedAt, time.Second)
}

// func TestDeleteMedia(t *testing.T) {
// 	media1 := createRandomMedia(t)
// 	err := testStore.DeleteMedia(context.Background(), media1.ID)
// 	require.NoError(t, err)

// 	media2, err := testStore.GetMedia(context.Background(), media1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, media2)
// }

func TestListMedia(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomMedia(t)
	}

	arg := ListMediaParams{
		Limit: 3,
		Offset: 3,
	}

	medium, err := testStore.ListMedia(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, medium, 3)

	for _, media := range medium {
		require.NotEmpty(t, media)
	}
}