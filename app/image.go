package app

import (
	"bytes"
	"image"
	"os"

	"fyne.io/fyne/v2"
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

func ToImage(data []byte) (image.Image, error) {
	buf := bytes.NewBuffer(data)
	img, _, err := image.Decode(buf)

	return img, err
}

func ToSize(size, max fyne.Size) fyne.Size {
	test := max.Width / size.Width
	height := size.Height * test

	if height > max.Height {
		test = max.Height / size.Height
	}

	return fyne.NewSize(size.Width*test, size.Height*test)
}
