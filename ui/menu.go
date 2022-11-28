package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func NewMenu() *fyne.MainMenu {
	open := fyne.NewMenuItem("Oepn", func() {
		log.Println("open")
	})
	open.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
	file := fyne.NewMenu("File", open)
	menus := fyne.NewMainMenu(file)

	return menus
}
