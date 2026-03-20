package image

import (
	"bytes"
	"errors"
	_ "github.com/gen2brain/heic"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
)

func ImageProccesor(photoBytes []byte, maxSide int) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(photoBytes))
	if err != nil {
		return nil, errors.New("apperrors decoding image: " + err.Error())
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if format == "jpeg" && (width <= maxSide && height <= maxSide) {
		return photoBytes, nil
	}

	var newWidth, newHeight int

	if width > height {
		if width > maxSide {
			newWidth = maxSide
			newHeight = (height * maxSide) / width
		} else {
			newWidth = width
			newHeight = height
		}
	} else {
		if height > maxSide {
			newHeight = maxSide
			newWidth = (width * maxSide) / height
		} else {
			newHeight = height
			newWidth = width
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, dst, &jpeg.Options{Quality: 100}); err != nil {
		return nil, errors.New("apperrors encoding image to JPEG: " + err.Error())
	}
	return buf.Bytes(), nil
}
