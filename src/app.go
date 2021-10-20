package godo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func Start() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("GoDo")
	w.Resize(fyne.NewSize(300, 500))
	w.SetContent(container.NewMax())
	w.ShowAndRun()
}
