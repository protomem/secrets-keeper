package cryptor

import (
	"errors"
)

var (
	ErrInvalidData      = errors.New("invalid data")
	ErrInvalidDataSize  = errors.New("invalid data size")
	ErrInvalidKeySize   = errors.New("invalid key size")
	ErrInvalidBlockSize = errors.New("invalid block size")
	ErrInvalidPadding   = errors.New("invalid padding")
)

type Encryptor interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(data []byte, key []byte) ([]byte, error)
}

type Encoder interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}

type Paddinger interface {
	Padding(data []byte, blockSize int) ([]byte, error)
	Unpadding(data []byte, blockSize int) ([]byte, error)
}
