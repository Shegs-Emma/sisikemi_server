package gapi

import (
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		ProfilePhoto: user.ProfilePhoto,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func convertMedia(media db.Medium) *pb.Media {
	return &pb.Media{
		MediaRef: media.MediaRef,
		Url: media.Url,
		AwsId: media.AwsID,
		CreatedAt: timestamppb.New(media.CreatedAt),
	}
}