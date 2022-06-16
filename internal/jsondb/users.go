package jsondb

import (
	"encoding/json"
	"errors"
	"gitlab.com/g6834/team41/auth/internal/models"
	"io"
	"os"
	"sync"
)

type JsonUsers struct {
	file os.File
	// fileMutex should be used when something accessing JsonUsers.file.
	fileMutex sync.Mutex
	// atomicityMutex should be used when we have more than one query associated with one request.
	atomicityMutex sync.Mutex
}

func NewJsonUsers(path string) (*JsonUsers, error) {
	// Open the file or create it if it doesn't exist.
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create(path)
			if err != nil {
				return nil, err
			}
			_, err = f.Write([]byte("[]"))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	db := &JsonUsers{
		file:           *f,
		fileMutex:      sync.Mutex{},
		atomicityMutex: sync.Mutex{},
	}

	return db, nil
}

func (db *JsonUsers) read() ([]models.User, error) {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	_, err := db.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(&db.file)
	if err != nil {
		return nil, err
	}

	var res []models.User
	err = json.Unmarshal(bytes, &res)

	return res, err
}

func (db *JsonUsers) write(m []models.User) error {
	db.fileMutex.Lock()
	defer db.fileMutex.Unlock()

	_, err := db.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = db.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = db.file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

// ChangeToken changing token of the user with given login.
func (db *JsonUsers) ChangeToken(token, login string) error {
	db.atomicityMutex.Lock()
	defer db.atomicityMutex.Unlock()

	users, err := db.read()
	if err != nil {
		return err
	}

	for i := range users {
		if users[i].Login == login {
			users[i].Token = token
		}
	}

	return db.write(users)
}

// GetUser selecting one user by login
func (db *JsonUsers) GetUser(login string) (*models.User, error) {
	users, err := db.read()
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].Login == login {
			return &users[i], nil
		}
	}

	return nil, errors.New("user isn't exist")
}

func (db *JsonUsers) AddUser(user models.User) error {
	db.atomicityMutex.Lock()
	defer db.atomicityMutex.Unlock()

	users, err := db.read()
	if err != nil {
		return err
	}

	users = append(users, user)
	return db.write(users)
}
