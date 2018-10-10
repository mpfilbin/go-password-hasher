package password

import (
	"crypto/sha512"
	"encoding/base64"
)

func Encode(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
