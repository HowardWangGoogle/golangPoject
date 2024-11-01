package gapi

import (
	"fmt"

	"github.com/techschool/simplebank/dbsq"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config          util.Config
	store           dbsq.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store dbsq.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmeticricKEy)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
