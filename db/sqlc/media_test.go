package db

import (
	"context"
	"testing"
	"time"

	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/stretchr/testify/require"
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