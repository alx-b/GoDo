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

type TaskDB struct {
	data       *sql.DB
	driverName string
	filePath   string
}

func CreateDB(driverName, filePath string) *TaskDB {
	db, err := sql.Open(driverName, filePath)
	if err != nil {
		return nil
	}
	createIfNotExist(db)
	return &TaskDB{db, driverName, filePath}
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

func (db *TaskDB) GetFromDB() ([]Task, error) {
	rows, err := db.data.Query("SELECT name, status FROM task")

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

func (db *TaskDB) AddToDB(task Task) error {
	_, err := db.data.Exec("INSERT INTO task (name, status) VALUES (?,?)",
		task.name, task.status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *TaskDB) UpdateToDB(oldName string, task Task) error {
	_, err := db.data.Exec("Update task SET name=?,status=? WHERE name=?", task.name, task.status, oldName)

	if err != nil {
		return err
	}

	return nil
}

func (db *TaskDB) DeleteFromDB(name string) error {
	_, err := db.data.Exec("DELETE FROM task WHERE name=?", name)

	if err != nil {
		return err
	}

	return nil
}

func (db *TaskDB) DeleteAllFromDB() error {
	_, err := db.data.Exec("DELETE FROM task")

	if err != nil {
		return err
	}

	return nil
}

func (db *TaskDB) CloseConnection() {
	db.data.Close()
}
