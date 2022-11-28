package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func NewImage() *canvas.Image {
	img := canvas.NewImageFromResource(theme.FyneLogo())
	// img.FillMode = canvas.ImageFillContain
	img.FillMode = canvas.ImageFillOriginal

	return img
}
