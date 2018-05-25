package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(orig string) string {
	return strings.ToLower(hex.EncodeToString(_md5([]byte(orig))))
}

func _md5(orig []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(orig)
	cipherStr := md5Ctx.Sum(nil)
	// fmt.Print(cipherStr)
	// fmt.Print("\n")
	// fmt.Print()
	return cipherStr
}
