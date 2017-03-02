package notorm

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_dbDriver = "mysql"
	_dbParam  = "root:rootroot@tcp(127.0.0.1:3306)/notorm"
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

type Article struct {
	ID      int
	Content string `mysql:"type=TEXT"`
}

func TestAll(t *testing.T) {
	no := NewConnection(_dbDriver, _dbParam)
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
	no := NewConnection(_dbDriver, _dbParam)
	no.Debug(true)
	no.CreateTable(User{})
	no.CreateTable(Email{})
	no.Insert(User{2, "Simon", "OCT"})
	no.Insert(Email{3, "simon3@foo.com", 2})
	no.Insert(Email{4, "simon4@foo.com", 2})
	no.Insert(Email{5, "simon5@foo.com", 2})
	no.Insert(Email{6, "simon6@foo.com", 2})
	arr, err := no.SelectAll("WHERE userid=2", Email{})
	if err != nil {
		t.Errorf("failed.")
	}
	if len(arr) != 4 {
		t.Errorf("should have 4 items: %d", len(arr))
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

func TestType(t *testing.T) {
	no := NewConnection(_dbDriver, _dbParam)
	no.Debug(true)
	no.CreateTable(Article{})
}
