package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func NewMain(
	app fyne.App,
	img *canvas.Image,
	toolbar *widget.Toolbar,
	menu *fyne.MainMenu,
) fyne.Window {
	main := app.NewWindow("Comic")
	main.Resize(fyne.NewSize(500, 500))
	img.SetMinSize(fyne.NewSize(300, 300))
	scroll := container.NewScroll(img)
	scroll.Resize(fyne.NewSize(600, 600))
	content := container.NewBorder(toolbar, nil, nil, nil, scroll)

	main.SetContent(content)
	main.SetMainMenu(menu)
	// main.SetFullScreen(true)
	plus := &desktop.CustomShortcut{KeyName: fyne.KeyPlus, Modifier: fyne.KeyModifierControl}
	w := 600
	main.Canvas().AddShortcut(plus, func(shortcut fyne.Shortcut) {
		log.Println("+", img.Size(), w)
		w += 100
		img.FillMode = canvas.ImageFillStretch
		img.Resize(fyne.NewSize(float32(w), float32(w)))
		img.FillMode = canvas.ImageFillOriginal
	})
	equal := &desktop.CustomShortcut{KeyName: fyne.KeyEqual, Modifier: fyne.KeyModifierControl}
	main.Canvas().AddShortcut(equal, func(shortcut fyne.Shortcut) {
		log.Println("+")
	})

	return main
}
