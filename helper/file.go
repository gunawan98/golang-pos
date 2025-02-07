package helper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveFile(file multipart.File, filename string) string {
	dir := "uploads/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(dir, newFilename)
	dst, err := os.Create(filePath)
	PanicIfError(err)
	defer dst.Close()

	_, err = io.Copy(dst, file)
	PanicIfError(err)

	return filePath
}
