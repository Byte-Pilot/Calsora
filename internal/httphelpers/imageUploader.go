package httphelpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
)

func UploadImage(c *gin.Context, maxSize int64) ([]byte, error) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return nil, errors.New("no file")
	}

	if fileHeader.Size > maxSize {
		return nil, errors.New("file too big")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("failed to read file")
	}

	return fileBytes, nil
}
