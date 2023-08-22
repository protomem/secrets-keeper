package model

import (
	"errors"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	ID int `json:"id"`

	CreatedAt time.Time `json:"createdAt"`

	Message string `json:"message"`
}
