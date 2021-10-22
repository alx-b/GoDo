package godo

import (
	"reflect"
	"testing"
)

func TestDatabaseQueries(t *testing.T) {
	db := CreateDB("sqlite3", ":memory:")

	t.Run("Create and connect to database", func(t *testing.T) {
		db := CreateDB("sqlite3", ":memory:")

		if db == nil {
			t.Errorf("Couldn't connect to database")
		}
	})

	t.Run("Add a task to database", func(t *testing.T) {
		db.AddToDB(Task{"testing1", ToDo})
		db.AddToDB(Task{"testing2", InProgress})
		got, err := db.GetFromDB()
		want := []Task{{"testing1", ToDo}, {"testing2", InProgress}}
		assertDeepEqual(t, got, want)
		assertErrors(t, err)
	})

	t.Run("Update a task to database", func(t *testing.T) {
		db.UpdateToDB("testing1", Task{"updatedName", ToDo})
		got, err := db.GetFromDB()
		want := []Task{{"updatedName", ToDo}, {"testing2", InProgress}}
		assertDeepEqual(t, got, want)
		assertErrors(t, err)
	})

	t.Run("Delete a task from database", func(t *testing.T) {
		db.DeleteFromDB("updatedName")
		got, err := db.GetFromDB()
		want := []Task{{"testing2", InProgress}}
		assertDeepEqual(t, got, want)
		assertErrors(t, err)
	})

	t.Run("Delete all task from database", func(t *testing.T) {
		db.DeleteAllFromDB()
		got, err := db.GetFromDB()
		want := []Task{}
		assertDeepEqual(t, got, want)
		assertErrors(t, err)
	})

}

func assertDeepEqual(t testing.TB, got, want []Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertErrors(t testing.TB, err error) {
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

/*
func TestGetCSVFile(t *testing.T) {
	t.Run(
		"Open csv file, only check if it opens a file not necessarely the same file",
		func(t *testing.T) {
			testCSVFile, err := os.Open("test.csv")
			if err != nil {
				testCSVFile, err = os.Create("test.csv")
			}
			got := GetCSVFile()
			want := testCSVFile

			if reflect.DeepEqual(got, want) {
				t.Errorf("got %p want %p", got, want)
			}
		},
	)
}

func TestWriteToFile(t *testing.T) {
	t.Run("Write Task to csv file", func(t *testing.T) {
		WriteToFile(GetCSVFile(), Task{"task1", 0})
		//want := []string{"task1", "0"}

	})
}
*/
