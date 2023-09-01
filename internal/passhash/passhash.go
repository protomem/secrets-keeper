package passhash

import "errors"

var ErrWrongPassword = errors.New("wrong password")

type Hasher interface {
	Generate(password string) (string, error)
	Compare(password string, hash string) error
}
