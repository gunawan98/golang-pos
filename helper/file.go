package helper

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// generateShortUniqueName generates a short unique name using timestamp and random bytes
func generateShortUniqueName() (string, error) {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes) + "_" + strconv.FormatInt(timestamp, 36), nil
}

func SaveFile(file multipart.File, filename string) string {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	PanicIfError(err)

	// Generate a short unique name
	uniqueName, err := generateShortUniqueName()
	PanicIfError(err)

	uploadParams := uploader.UploadParams{
		Folder:   "product",
		PublicID: uniqueName,
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploadParams)
	PanicIfError(err)

	return uploadResult.SecureURL
}

func DeleteFile(url string) error {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		fmt.Println("Error creating Cloudinary instance:", err)
		return err
	}

	// Extract the public ID from the Cloudinary URL
	publicID := extractPublicID(url)

	// Delete the image from Cloudinary using the Admin API
	_, err = cld.Admin.DeleteAssets(context.Background(), admin.DeleteAssetsParams{
		PublicIDs:    []string{"product/" + publicID},
		DeliveryType: "upload",
		AssetType:    "image",
	})
	if err != nil {
		return err
	}

	return nil
}

// extractPublicID extracts the public ID from the Cloudinary URL
func extractPublicID(url string) string {
	// Assuming the URL format is something like: https://res.cloudinary.com/<cloud_name>/image/upload/v<version>/<public_id>.<format>
	parts := strings.Split(url, "/")
	publicIDWithFormat := parts[len(parts)-1]
	publicID := strings.Split(publicIDWithFormat, ".")[0]
	return publicID
}
