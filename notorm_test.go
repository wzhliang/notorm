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
	err := no.SelectAll("WHERE id=1", &u)
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
	err = no.SelectAll("WHERE id=3", &e)
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
