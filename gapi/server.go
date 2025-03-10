package gapi

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSisikemiFashionServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
	taskdistributor worker.TaskDistributor
	cloud *cloudinary.Cloudinary
}

func NewServer (config util.Config, store db.Store, taskDistributor worker.TaskDistributor, cloud *cloudinary.Cloudinary) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server {
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskdistributor: taskDistributor,
		cloud: cloud,
	}

	return server, nil
}