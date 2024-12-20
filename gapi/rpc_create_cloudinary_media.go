package gapi

import (
	"bytes"
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/techschool/simplebank/pb"
)

func (server *Server) UploadCloudinaryMedia (ctx context.Context, req *pb.UploadCloudinaryImageRequest) (*pb.UploadCloudinaryImageResponse, error) {
	uploadResult, err := server.cloud.Upload.Upload(ctx, bytes.NewReader(req.ImageData), uploader.UploadParams{
		PublicID: req.ImageName,
	})

	if err != nil {
		return nil, err
	}

	return &pb.UploadCloudinaryImageResponse{
		Url: uploadResult.SecureURL,
	}, nil
}