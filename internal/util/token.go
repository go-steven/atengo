package util

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func Token() string {
	token := uuid.NewV4()
	return strings.ToUpper(Sha1(token.String())[:aes.BlockSize*2])
}

func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
