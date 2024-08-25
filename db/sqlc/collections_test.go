package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomCollection(t *testing.T) Collection {
	collectionName := util.RandomString(6)

	collection, err := testQueries.CreateCollection(context.Background(), collectionName)

	require.NoError(t, err)
	require.NotEmpty(t, collection)
	require.Equal(t, collectionName, collection.CollectionName)
	require.NotZero(t, collection.ID)
	require.NotZero(t, collection.CreatedAt)

	return collection
}

func TestCreateCollection(t *testing.T) {
	createRandomCollection(t)
}

func TestGetCollection(t *testing.T) {
	collection1 := createRandomCollection(t)
	collection2, err := testQueries.GetCollection(context.Background(), collection1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, collection2)
	require.Equal(t, collection1.ID, collection2.ID)
	require.Equal(t, collection1.CollectionName, collection2.CollectionName)
	require.Equal(t, collection1.LastUpdatedAt, collection2.LastUpdatedAt)
	require.WithinDuration(t, collection1.CreatedAt, collection2.CreatedAt, time.Second)
}

func TestListCollection(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomCollection(t)
	}

	arg := ListCollectionParams{
		Limit: 3,
		Offset: 3,
	}

	collections, err := testQueries.ListCollection(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, collections, 3)
	
	for _, collection := range collections {
		require.NotEmpty(t, collection)
	}
}