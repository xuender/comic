//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/xuender/comic/app"
)

func InitApp() *app.App {
	wire.Build(
		app.NewApp,
		app.NewCache,
		app.NewFiles,
	)
	return &app.App{}
}
