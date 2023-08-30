package model

import (
	"errors"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	ID int `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`

	AccessKey  string `json:"-"`
	SigningKey string `json:"-"`

	SecretPhrase string `json:"-"`

	Message string `json:"message"`
}
