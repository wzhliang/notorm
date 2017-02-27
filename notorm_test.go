package notorm

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID      int
	Name    string
	Address string
}

type Email struct {
	ID     int
	Email  string
	UserID int
}

func TestAll(t *testing.T) {
	no := NewConnection("sqlite3", "notorm.db")
	no.CreateTable(User{})
	no.CreateTable(Email{})
	no.Insert(User{1, "Simon", "OCT"})
	u := User{}
	err := no.Select("WHERE id=1", &u)
	if err != nil {
		t.Error("Failed to select")
	}
	if u.ID != 1 {
		t.Errorf("id error: %d\n", u.ID)
	}
	if u.Name != "Simon" {
		t.Errorf("Name error: %d\n", u.Name)
	}
	if u.Address != "OCT" {
		t.Errorf("address error: %d\n", u.Address)
	}
	no.Insert(Email{3, "simon@foo.com", 1})
	e := Email{}
	err = no.Select("WHERE id=3", &e)
	if err != nil {
		t.Error("Failed to select")
	}
	if e.ID != 3 {
		t.Errorf("id error: %d\n", e.ID)
	}
	if e.Email != "simon@foo.com" {
		t.Errorf("email error: %d\n", e.Email)
	}
	if e.UserID != 1 {
		t.Errorf("uid error: %d\n", e.UserID)
	}
}

func TestSelectAll(t *testing.T) {
	no := NewConnection("sqlite3", "notorm1.db")
	no.CreateTable(User{})
	no.CreateTable(Email{})
	no.Debug(true)
	no.Insert(User{1, "Simon", "OCT"})
	no.Insert(Email{3, "simon3@foo.com", 1})
	no.Insert(Email{4, "simon4@foo.com", 1})
	no.Insert(Email{5, "simon5@foo.com", 1})
	no.Insert(Email{6, "simon6@foo.com", 1})
	arr, err := no.SelectAll("WHERE userid=1", Email{})
	if err != nil {
		t.Errorf("failed.")
	}
	if len(arr) != 4 {
		t.Errorf("should have 4 items")
	}
	email := arr[0].(*Email)
	if email.ID != 3 {
		t.Errorf("Wrong email id: %d\n", email.ID)
	}
	email = arr[3].(*Email)
	if email.ID != 6 {
		t.Errorf("Wrong email id: %d\n", email.ID)
	}
}
