package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func Sha1(orig string) string {
	return strings.ToLower(hex.EncodeToString(Sha1Impl([]byte(orig))))
}

func Sha1Impl(orig []byte) []byte {
	ctx := sha1.New()
	ctx.Write(orig)
	cipherStr := ctx.Sum(nil)
	// fmt.Print(cipherStr)
	// fmt.Print("\n")
	// fmt.Print()
	return cipherStr
}
