package mongo

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.com/g6834/team41/auth/internal/models"
)

// TODO move to integration tests

const (
	// harder better faster stronger
	hardcoded_TEST_DB_PORT     = 27017
	hardcoded_TEST_DB_LOGIN    = "team41"
	hardcoded_TEST_DB_PASSWORD = "mNgdxQTbhVGd"
	hardcoded_TEST_DB_ADDRESS  = "91.185.93.34"
	hardcoded_TEST_DB_NAME     = "team41"
)

/* Copy-paste pool
port, err := strconv.Atoi(os.Getenv("TEST_DB_PORT"))
db, err = NewUsers(os.Getenv("TEST_DB_LOGIN"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_ADDRESS"), os.Getenv("TEST_DB_NAME"), port)
*/

func TestAddGet(t *testing.T) {
	var db *Users
	var err error
	t.Run("connecting to db", func(t *testing.T) {
		db, err = NewUsers(
			hardcoded_TEST_DB_LOGIN,
			hardcoded_TEST_DB_PASSWORD,
			hardcoded_TEST_DB_ADDRESS,
			hardcoded_TEST_DB_NAME,
			hardcoded_TEST_DB_PORT,
		)
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

	t.Run("change token", func(t *testing.T) {
		otherToken := "other-test"
		err := db.ChangeToken(otherToken, u.Login)
		if err != nil {
			t.Error(err)
		}

		upd, err := db.GetUser(u.Login)
		if err != nil {
			t.Error(err)
		}

		if upd.Token != otherToken {
			t.Errorf("wrong token, expected %v but got %v", otherToken, upd.Token)
		}
	})

	t.Cleanup(func() {
		err = db.DeleteUser(u.Login)
		if err != nil {
			t.Error(err)
		}
	})
}
