package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(orig string) string {
	return strings.ToLower(hex.EncodeToString(Md5Impl([]byte(orig))))
}

func Md5Impl(orig []byte) []byte {
	ctx := md5.New()
	ctx.Write(orig)
	cipherStr := ctx.Sum(nil)
	// fmt.Print(cipherStr)
	// fmt.Print("\n")
	// fmt.Print()
	return cipherStr
}
