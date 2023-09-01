package base64

import (
	"encoding/base64"
	"fmt"

	"github.com/protomem/secrets-keeper/internal/cryptor"
)

var _ cryptor.Encoder = (*Encoder)(nil)

type Encoder struct {
	isURL bool
}

func NewEncoder(isURL bool) *Encoder {
	return &Encoder{
		isURL: isURL,
	}
}

func (e *Encoder) Encode(data []byte) ([]byte, error) {
	const _ = "base64.Encode"

	encodedDataSize := e.encoding().EncodedLen(len(data))
	encodedData := make([]byte, encodedDataSize)
	e.encoding().Encode(encodedData, data)

	return encodedData, nil
}

func (e *Encoder) Decode(data []byte) ([]byte, error) {
	const op = "base64.Decode"

	decodedDataSize := e.encoding().DecodedLen(len(data))
	decodedData := make([]byte, decodedDataSize)
	_, err := e.encoding().Decode(decodedData, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return decodedData, nil
}

func (e *Encoder) encoding() *base64.Encoding {
	if e.isURL {
		return base64.RawURLEncoding
	} else {
		return base64.RawStdEncoding
	}
}
