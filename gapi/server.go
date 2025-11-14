package gapi

import (
	"fmt"

	db "github.com/Shegs-Emma/sisikemi_server/db/sqlc"
	"github.com/Shegs-Emma/sisikemi_server/pb"
	"github.com/Shegs-Emma/sisikemi_server/token"
	"github.com/Shegs-Emma/sisikemi_server/util"
	"github.com/Shegs-Emma/sisikemi_server/worker"
	"github.com/cloudinary/cloudinary-go/v2"
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