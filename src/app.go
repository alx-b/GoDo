package godo

import (
	"fmt"
	"reflect"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Start() {
	db := CreateDB("sqlite3", "godo_data.db")
	defer db.CloseConnection()

	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
	window := app.NewWindow("GoDo")
	window.Resize(fyne.NewSize(300, 500))
	window.SetContent(mainContent(app, window, db))
	window.ShowAndRun()
}

func splitTaskListByStatus(db *TaskDB) ([]string, []string, []string) {
	taskList, _ := db.GetFromDB()

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

func mainContent(a fyne.App, w fyne.Window, db *TaskDB) *fyne.Container {
	myListToDo, myListInProgress, myListDone := splitTaskListByStatus(db)

	listView := createListView(&myListToDo)
	listViewInProgress := createListView(&myListInProgress)
	listViewDone := createListView(&myListDone)

	listView.OnSelected = func(id widget.ListItemID) {
		db.UpdateToDB(myListToDo[id], Task{myListToDo[id], InProgress})
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(db)
		listView.UnselectAll()
		listView.Refresh()
	}

	listViewInProgress.OnSelected = func(id widget.ListItemID) {
		db.UpdateToDB(myListInProgress[id], Task{myListInProgress[id], Done})
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(db)
		listViewInProgress.UnselectAll()
		listViewInProgress.Refresh()
	}

	listViewDone.OnSelected = func(id widget.ListItemID) {
		db.DeleteFromDB(myListDone[id])
		myListToDo, myListInProgress, myListDone = splitTaskListByStatus(db)
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
			w.SetContent(addTaskContent(a, w, db))
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			db.DeleteAllFromDB()
			myListToDo, myListInProgress, myListDone = splitTaskListByStatus(db)
			tabs.Refresh()
		}),
	)

	return container.NewBorder(toolbar, nil, nil, nil, tabs)
}

func addTaskContent(a fyne.App, w fyne.Window, db *TaskDB) *fyne.Container {
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
			db.AddToDB(Task{taskEntry.Text, taskStat})
			taskEntry.Text = ""
			taskEntry.Refresh()
		}
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			w.SetContent(mainContent(a, w, db))
		}),
	)

	form := container.NewVBox(taskEntry, taskStatusSelection, submitButton)

	return container.NewBorder(toolbar, nil, nil, nil, form)
}
