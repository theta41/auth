package mongo

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/g6834/team41/auth/internal/models"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	var db *Users
	var err error
	t.Run("connecting to db", func(t *testing.T) {
		port, err := strconv.Atoi(os.Getenv("TEST_DB_PORT"))
		if err != nil {
			t.Fatal(err)
		}
		db, err = NewUsers(os.Getenv("TEST_DB_LOGIN"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_ADDRESS"), os.Getenv("TEST_DB_NAME"), port)
		if err != nil {
			t.Error(fmt.Errorf("error creating db: %w", err))
		}
	})

	u := models.User{
		Login:        uuid.New().String(),
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
		err = db.DeleteUser(u.Login)
		if err != nil {
			t.Error(err)
		}
	})
}
