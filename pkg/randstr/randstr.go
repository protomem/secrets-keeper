package randstr

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	_letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_letterIdxBits = 6
	_letterIdxMask = 1<<_letterIdxBits - 1
	_letterIdxMax  = 63 / _letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func Gen(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), _letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), _letterIdxMax
		}
		if idx := int(cache & _letterIdxMask); idx < len(_letterBytes) {
			b[i] = _letterBytes[idx]
			i--
		}
		cache >>= _letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
