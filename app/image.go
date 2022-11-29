package app

import (
	"bytes"
	"image"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/h2non/filetype"
)

func IsImage(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	defer file.Close()

	head := make([]byte, _headSize)
	_, _ = file.Read(head)

	return filetype.IsImage(head)
}

func IsArchive(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	defer file.Close()

	head := make([]byte, _headSize)
	_, _ = file.Read(head)

	return filetype.IsArchive(head)
}

func Image(path string) *canvas.Image {
	image, _ := ReadImage(path)
	log.Println(image.Bounds())
	img := canvas.NewImageFromImage(image)
	img.FillMode = canvas.ImageFillStretch
	img.ScaleMode = canvas.ImageScaleFastest

	return img
}

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
