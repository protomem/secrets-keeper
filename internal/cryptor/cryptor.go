package cryptor

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	const op = "cryptor.Encrypt"
	var err error

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	src, err := padding(data, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dst := make([]byte, len(src))

	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(dst, src)

	return dst, nil
}

func Decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	const op = "crytor.Decrypt"
	var err error

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dst := make([]byte, len(encryptedData))

	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(dst, encryptedData)

	dst, err = unpadding(dst, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dst, nil
}

func Encode(data []byte) (string, error) {
	encData := hex.EncodeToString(data)
	return encData, nil
}

func Decode(data string) ([]byte, error) {
	const op = "cryptor.Decode"
	var err error

	rawData := []byte(data)
	decData := make([]byte, hex.DecodedLen(len(rawData)))
	_, err = hex.Decode(decData, rawData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return decData, nil
}

func padding(src []byte, blockSize int) ([]byte, error) {
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("pkcs7: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(src)%blockSize
		padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(src, padding...), nil
	}
}

func unpadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	}
	padLen := int(src[length-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(src, ref) {
		log.Printf("padLen: %d blockSize: %d", padLen, blockSize)
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return src[:length-padLen], nil
}
