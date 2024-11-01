package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/dbsq"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store dbsq.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmeticricKEy: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
