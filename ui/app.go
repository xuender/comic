package ui

import "fyne.io/fyne/v2"

type App struct {
	app  fyne.App
	main fyne.Window
}

func NewApp(
	app fyne.App,
	main fyne.Window,
) *App {
	return &App{
		app:  app,
		main: main,
	}
}

func (p *App) Run() {
	p.main.Show()
	p.app.Run()
}
