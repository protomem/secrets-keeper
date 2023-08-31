package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/protomem/secrets-keeper/internal/cryptor"
)

var _ cryptor.Encryptor = (*Encryptor)(nil)

type Encryptor struct {
	encoder   cryptor.Encoder
	paddinger cryptor.Paddinger
}

func NewEncryptor(encoder cryptor.Encoder, paddinger cryptor.Paddinger) *Encryptor {
	return &Encryptor{
		encoder:   encoder,
		paddinger: paddinger,
	}
}

func (e *Encryptor) Encrypt(data []byte, key []byte) ([]byte, error) {
	const op = "aes.Encrypt"
	var err error

	block, err := aes.NewCipher(key)
	if err != nil {
		var aesKeySizeError aes.KeySizeError
		if errors.As(err, &aesKeySizeError) {
			return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidKeySize)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	alignedData, err := e.paddinger.Padding(data, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ecryptedData := make([]byte, len(alignedData))
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(ecryptedData, alignedData)

	encodedData, err := e.encoder.Encode(ecryptedData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return encodedData, nil
}

func (e *Encryptor) Decrypt(data []byte, key []byte) ([]byte, error) {
	const op = "aes.Decrypt"
	var err error

	decodedData, err := e.encoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		var aesKeySizeError aes.KeySizeError
		if errors.As(err, &aesKeySizeError) {
			return nil, fmt.Errorf("%s: %w", op, cryptor.ErrInvalidKeySize)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	decryptedData := make([]byte, len(decodedData))
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(decryptedData, decodedData)

	originData, err := e.paddinger.Unpadding(decryptedData, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return originData, nil
}
