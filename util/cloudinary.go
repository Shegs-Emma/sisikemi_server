package util

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary(config Config) *cloudinary.Cloudinary {
	cloud, err := cloudinary.NewFromURL(config.CloudinaryUrl)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	return cloud
}