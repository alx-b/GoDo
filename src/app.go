package godo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"reflect"
)

func Start() {
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow("GoDo")
	window.Resize(fyne.NewSize(300, 500))
	window.SetContent(mainContent(app, window))
	window.ShowAndRun()
}

func mainContent(a fyne.App, w fyne.Window) *fyne.Container {
	tabs := container.NewAppTabs(
		container.NewTabItem("To Do", widget.NewLabel("test1")),
		container.NewTabItem("In Progress", widget.NewLabel("test2")),
		container.NewTabItem("Done", widget.NewLabel("test3")),
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ColorChromaticIcon(), func() {
			if reflect.DeepEqual(a.Settings().Theme(), theme.DarkTheme()) {
				a.Settings().SetTheme(theme.LightTheme())
			} else {
				a.Settings().SetTheme(theme.DarkTheme())
			}
		}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			w.SetContent(addTaskContent(a, w))
		}),
	)

	return container.NewBorder(toolbar, nil, nil, nil, tabs)
}

func addTaskContent(a fyne.App, w fyne.Window) *fyne.Container {
	taskEntry := widget.NewEntry()
	taskEntry.SetPlaceHolder("Enter task name")
	taskStatusSelection := widget.NewSelect(
		[]string{"ToDo", "Doing", "Done"}, func(change string) {
		})
	taskStatusSelection.PlaceHolder = "ToDo"

	submitButton := widget.NewButton("Add", func() {
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			w.SetContent(mainContent(a, w))
		}),
	)

	form := container.NewVBox(taskEntry, taskStatusSelection, submitButton)

	return container.NewBorder(toolbar, nil, nil, nil, form)
}
