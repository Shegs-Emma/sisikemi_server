package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/token"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/Shegs-Emma/sisikemi_server/worker"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor, cloud *cloudinary.Cloudinary) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor, cloud)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, isAdmin bool, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration, isAdmin)
	require.NoError(t, err)
	
	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}
	
	return metadata.NewIncomingContext(context.Background(), md)
}