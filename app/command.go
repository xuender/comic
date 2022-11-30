package app

import "fyne.io/fyne/v2"

type Command struct {
	Help string
	Key1 fyne.KeyName
	Key2 fyne.KeyName
	Call func()
	Icon fyne.Resource
}
