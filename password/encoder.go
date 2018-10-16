package password

import (
	"crypto/sha512"
	"encoding/base64"
)

//Encode returns a base 64 encoded 512 bit SHA digest of a password string
func Encode(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
