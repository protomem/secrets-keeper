package model

import (
	"errors"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	ID int `json:"id"`

	CreatedAt time.Time `json:"createdAt"`

	AccessKey  string `json:"-"`
	SigningKey string `json:"-"`

	Message string `json:"message"`
}
