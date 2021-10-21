package godo

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"reflect"
	"sort"
)

func Start() {
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow("GoDo")
	window.Resize(fyne.NewSize(300, 500))
	window.SetContent(mainContent(app, window))
	window.ShowAndRun()
}

func splitTaskListByStatus(taskList []Task) ([]string, []string, []string) {
	myListToDo := []string{}
	myListInProgress := []string{}
	myListDone := []string{}

	for _, item := range taskList {
		switch item.status {
		case ToDo:
			myListToDo = append(myListToDo, item.name)
		case InProgress:
			myListInProgress = append(myListInProgress, item.name)
		case Done:
			myListDone = append(myListDone, item.name)
		}
	}

	sort.Strings(myListToDo)
	sort.Strings(myListInProgress)
	sort.Strings(myListDone)

	return myListToDo, myListInProgress, myListDone
}

func createListView(list *[]string) *widget.List {
	return widget.NewList(func() int {
		return len(*list)
	}, func() fyne.CanvasObject {
		return widget.NewLabel(".")
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		l := *list
		item.(*widget.Label).SetText(l[id])
	})
}

func mainContent(a fyne.App, w fyne.Window) *fyne.Container {
	list, _ := GetFromDB()
	myListToDo, myListInProgress, myListDone := splitTaskListByStatus(list)

	listView := createListView(&myListToDo)
	listViewInProgress := createListView(&myListInProgress)
	listViewDone := createListView(&myListDone)

	listView.OnSelected = func(id widget.ListItemID) {
		UpdateToDB(myListToDo[id], Task{myListToDo[id], InProgress})
		list, _ = GetFromDB()
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(list)
		listView.UnselectAll()
		listView.Refresh()
	}

	listViewInProgress.OnSelected = func(id widget.ListItemID) {
		UpdateToDB(myListInProgress[id], Task{myListInProgress[id], Done})
		list, _ = GetFromDB()
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(list)
		listViewInProgress.UnselectAll()
		listViewInProgress.Refresh()
	}

	listViewDone.OnSelected = func(id widget.ListItemID) {
		DeleteFromDB(myListDone[id])
		list, _ = GetFromDB()
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(list)
		listViewDone.UnselectAll()
		listViewDone.Refresh()
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("To Do", listView),
		container.NewTabItem("In Progress", listViewInProgress),
		container.NewTabItem("Done", listViewDone),
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
		[]string{"ToDo", "InProgress", "Done"}, func(change string) {
		})
	taskStatusSelection.PlaceHolder = "ToDo"

	submitButton := widget.NewButton("Add", func() {
		if taskEntry.Text == "" {
			fmt.Println("Need to write a task")
		} else {
			var taskStat Status
			switch taskStatusSelection.Selected {
			case "ToDo":
				taskStat = ToDo
			case "InProgress":
				taskStat = InProgress
			case "Done":
				taskStat = Done
			}
			AddToDB(Task{taskEntry.Text, taskStat})
			taskEntry.Text = ""
			taskEntry.Refresh()
		}
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			w.SetContent(mainContent(a, w))
		}),
	)

	form := container.NewVBox(taskEntry, taskStatusSelection, submitButton)

	return container.NewBorder(toolbar, nil, nil, nil, form)
}
