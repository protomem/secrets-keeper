package argon2

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"

	"github.com/protomem/secrets-keeper/internal/cryptor"
	"github.com/protomem/secrets-keeper/internal/passhash"
	"github.com/protomem/secrets-keeper/pkg/randstr"
	"golang.org/x/crypto/argon2"
)

var _ passhash.Hasher = (*Hasher)(nil)

var DefaultOptions = Options{
	Memory:     64 * 1024,
	Iterations: 3,
	Parallel:   2,
	SaltLength: 16,
	KeyLength:  32,
}

type Options struct {
	Memory     uint32
	Iterations uint32
	Parallel   uint8
	SaltLength uint32
	KeyLength  uint32
}

type Hasher struct {
	encoder cryptor.Encoder
	opts    Options
}

func NewHasher(encoder cryptor.Encoder, opts Options) *Hasher {
	return &Hasher{
		encoder: encoder,
		opts:    opts,
	}
}

func (h *Hasher) Generate(password string) (string, error) {
	const op = "argon2.Generate"

	salt := []byte(randstr.Gen(int(h.opts.SaltLength)))
	hash := argon2.IDKey(
		[]byte(password), salt,
		h.opts.Iterations, h.opts.Memory, h.opts.Parallel, h.opts.KeyLength,
	)

	encodedHash, err := h.encode(hash, salt)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return encodedHash, nil
}

func (h *Hasher) Compare(password string, hash string) error {
	const op = "argon2.Compare"

	decodedHash, decodedSalt, err := h.decode(hash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newHash := argon2.IDKey(
		[]byte(password), []byte(decodedSalt),
		h.opts.Iterations, h.opts.Memory, h.opts.Parallel, h.opts.KeyLength,
	)

	if subtle.ConstantTimeCompare(newHash, []byte(decodedHash)) == 1 {
		return nil
	}

	return fmt.Errorf("%s: %w", op, passhash.ErrWrongPassword)
}

func (h *Hasher) encode(hash []byte, salt []byte) (string, error) {
	const op = "encode"
	var err error

	encodedSalt, err := h.encoder.Encode(salt)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	encodedHash, err := h.encoder.Encode(hash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.opts.Memory, h.opts.Iterations, h.opts.Parallel,
		encodedSalt, encodedHash,
	), nil
}

func (h *Hasher) decode(encodedHash string) ([]byte, []byte, error) {
	const op = "decode"
	var err error

	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, fmt.Errorf("%s: %w", op, errors.New("invalid hash"))
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	if version != argon2.Version {
		return nil, nil, fmt.Errorf("%s: %w", op, errors.New("incorrect version"))
	}

	var (
		memory     uint32
		iterations uint32
		parallel   uint8
	)
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallel)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	if memory != h.opts.Memory || iterations != h.opts.Iterations || parallel != h.opts.Parallel {
		return nil, nil, fmt.Errorf("%s: %w", op, errors.New("incorrect options"))
	}

	salt, err := h.encoder.Decode([]byte(vals[4]))
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	hash, err := h.encoder.Decode([]byte(vals[5]))
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return hash, salt, nil
}
