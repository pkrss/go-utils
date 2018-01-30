package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(orig string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(orig))
	cipherStr := md5Ctx.Sum(nil)
	// fmt.Print(cipherStr)
	// fmt.Print("\n")
	// fmt.Print()
	return strings.ToLower(hex.EncodeToString(cipherStr))
}
