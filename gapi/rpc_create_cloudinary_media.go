package gapi

import (
	"bytes"
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UploadCloudinaryMedia (ctx context.Context, req *pb.UploadCloudinaryImageRequest) (*pb.UploadCloudinaryImageResponse, error) {
	uploadResult, err := server.cloud.Upload.Upload(ctx, bytes.NewReader(req.ImageData), uploader.UploadParams{
		PublicID: req.ImageName,
		Folder: "sisikemi_fashion/",
	})

	if err != nil {
		return nil, err
	}

	arg := db.CreateMediaParams{
		MediaRef: util.RandomString(10),
		Url: uploadResult.SecureURL,
		AwsID: uploadResult.AssetID,
	}

	result, err := server.store.CreateMedia(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	return &pb.UploadCloudinaryImageResponse{
		Media: convertMedia(result),
	}, nil
}