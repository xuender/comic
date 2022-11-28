//go:build wireinject
// +build wireinject

package cmd

import (
	"fyne.io/fyne/v2/app"
	"github.com/google/wire"
	"github.com/xuender/comic/ui"
)

func InitApp() *ui.App {
	wire.Build(
		app.New,
		ui.NewApp,
		ui.NewMain,
		ui.NewImage,
		ui.NewToolbar,
		ui.NewMenu,
	)
	return &ui.App{}
}
