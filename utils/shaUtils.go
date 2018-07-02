package utils

import (
	"crypto/sha256"
	"strings"
	"crypto/sha1"
	"crypto/sha512"
)

type SHAUtils struct {
}

func (self SHAUtils) Hash(data *string, shaType *string) ([]byte) {
	var h = sha256.New()
	if (strings.ToUpper(*shaType) == "SHA1") {
		h = sha1.New()
	} else if (strings.ToUpper(*shaType) == "SHA512") {
		h = sha512.New()
	}
	h.Write([]byte(*data))
	return h.Sum(nil);
}
