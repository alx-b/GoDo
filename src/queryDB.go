package godo

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Status int

const (
	ToDo Status = iota
	InProgress
	Done
)

type Task struct {
	name   string
	status Status
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "godo_data.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createIfNotExist(db *sql.DB) error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY, name TEXT, status INTEGER)",
	)

	if err != nil {
		return err
	}

	return nil
}

func GetFromDB() ([]Task, error) {
	db, err := openDB()
	defer db.Close()

	createIfNotExist(db)

	rows, err := db.Query("SELECT name, status FROM task")

	if err != nil {
		return nil, err
	}

	list := []Task{}
	for rows.Next() {
		myTask := Task{}
		rows.Scan(&myTask.name, &myTask.status)
		list = append(list, myTask)
	}

	return list, nil
}

func AddToDB(task Task) error {
	db, err := openDB()
	defer db.Close()

	_, err = db.Exec("INSERT INTO task (name, status) VALUES (?,?)",
		task.name, task.status,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateToDB(oldName string, task Task) error {
	db, err := openDB()
	defer db.Close()

	_, err = db.Exec("Update task SET name=?,status=? WHERE name=?", task.name, task.status, oldName)

	if err != nil {
		return err
	}

	return nil
}

func DeleteFromDB(name string) error {
	db, err := openDB()
	defer db.Close()

	_, err = db.Exec("DELETE FROM task WHERE name=?", name)

	if err != nil {
		return err
	}

	return nil
}

func DeleteAllFromDB() error {
	db, err := openDB()
	defer db.Close()

	_, err = db.Exec("DELETE FROM task")

	if err != nil {
		return err
	}

	return nil
}
