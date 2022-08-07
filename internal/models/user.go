package models

import "time"

type User struct {
	Login        string    `json:"login"`
	Token        string    `json:"token"`
	PasswordHash string    `json:"password_hash"`
	RegisterDate time.Time `json:"register_date"`
}

type CtxUsername struct{}
