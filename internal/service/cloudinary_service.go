package service

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/satyam-svg/resume-parser/config"
)

func UploadToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	cld, err := cloudinary.NewFromParams(
		config.AppConfig.CloudinaryCloudName,
		config.AppConfig.CloudinaryAPIKey,
		config.AppConfig.CloudinaryAPISecret,
	)
	if err != nil {
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   "profile_images",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
