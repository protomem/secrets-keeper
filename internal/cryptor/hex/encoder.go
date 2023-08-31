package hex

import (
	"encoding/hex"
	"fmt"

	"github.com/protomem/secrets-keeper/internal/cryptor"
)

var _ cryptor.Encoder = (*Encoder)(nil)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(data []byte) ([]byte, error) {
	const _ = "hex.Encode"

	encodedDataSize := hex.EncodedLen(len(data))
	encodedData := make([]byte, encodedDataSize)
	hex.Encode(encodedData, data)

	return encodedData, nil
}

func (e *Encoder) Decode(data []byte) ([]byte, error) {
	const op = "hex.Decode"

	decodedDataSize := hex.DecodedLen(len(data))
	decodedData := make([]byte, decodedDataSize)
	_, err := hex.Decode(decodedData, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return decodedData, nil
}
