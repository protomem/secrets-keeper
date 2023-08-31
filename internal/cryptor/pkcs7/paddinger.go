package pkcs7

import (
	"bytes"
	"fmt"

	"github.com/protomem/secrets-keeper/internal/cryptor"
)

var _ cryptor.Paddinger = (*Paddinger)(nil)

type Paddinger struct{}

func NewPaddinger() *Paddinger {
	return &Paddinger{}
}

func (p *Paddinger) Padding(data []byte, blockSize int) ([]byte, error) {
	const op = "pkcs7.Padding"

	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidBlockSize)
	}

	paddingSize := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	return append(data, padding...), nil
}

func (p *Paddinger) Unpadding(data []byte, blockSize int) ([]byte, error) {
	const op = "pkcs7.Unpadding"

	dataSize := len(data)
	if dataSize == 0 {
		return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidDataSize)
	}

	if dataSize%blockSize != 0 {
		return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidData)
	}

	paddingSize := int(data[dataSize-1])
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	if paddingSize > blockSize || paddingSize == 0 || !bytes.HasSuffix(data, padding) {
		return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidData)
	}

	return data[:dataSize-paddingSize], nil
}
