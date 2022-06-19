package jsondb

import (
	"gitlab.com/g6834/team41/auth/internal/models"
	"os"
	"reflect"
	"testing"
	"time"
)

const DBName = "test.json"

func TestAddGet(t *testing.T) {
	var err error
	db := &JsonUsers{}
	t.Run("create empty db", func(t *testing.T) {
		db, err = NewJsonUsers(DBName)
		if err != nil {
			t.Error(err)
		}
	})

	u := models.User{
		Login:        "test",
		Token:        "test",
		PasswordHash: "test",
		RegisterDate: time.Now(),
	}

	t.Run("insert into db", func(t *testing.T) {
		err := db.AddUser(u)
		if err != nil {
			t.Error(err)
		}
	})

	got := &models.User{}
	t.Run("select from db", func(t *testing.T) {
		got, err = db.GetUser(u.Login)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("compare results", func(t *testing.T) {
		if reflect.DeepEqual(*got, u) {
			t.Errorf("%v is not equal to %v\n", *got, u)
		}
	})

	t.Cleanup(func() {
		err = os.Remove(DBName)
		if err != nil {
			t.Error(err)
		}
	})
}
