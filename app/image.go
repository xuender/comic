package app

import (
	"image"
	"os"
)

func ReadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	return img, err
}
