package uuid

import (
	crand "crypto/rand"
	"encoding/hex"
	"strings"
)

func UuidRemoveLine(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.Replace(s, "-", "", -1)
}

// seeded indicates if math/rand has been seeded
// var seeded bool = false

// func randBytes(x []byte) {

// 	length := len(x)
// 	n, err := crand.Read(x)

// 	if n != length || err != nil {
// 		if !seeded {
// 			mrand.Seed(time.Now().UnixNano())
// 		}

// 		for length > 0 {
// 			length--
// 			x[length] = byte(mrand.Int31n(256))
// 		}
// 	}
// }

func setVersion(u *[16]byte, v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits.
func setVariant(u *[16]byte, v byte) {
	switch v {
	case 0:
		u[8] = (u[8]&(0xff>>1) | (0x00 << 7))
	case 1:
		u[8] = (u[8]&(0xff>>2) | (0x02 << 6))
	case 2:
		u[8] = (u[8]&(0xff>>3) | (0x06 << 5))
	case 3:
		fallthrough
	default:
		u[8] = (u[8]&(0xff>>3) | (0x07 << 5))
	}
}

func uuid2String(u *[16]byte) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func UuidCreate(keepLine ...bool) string {
	// var x [16]byte
	// randBytes(x[:])
	// x[6] = (x[6] & 0x0F) | 0x40
	// x[8] = (x[8] & 0x3F) | 0x80
	// return x

	u := [16]byte{}
	if _, err := crand.Reader.Read(u[:]); err != nil {
		return ""
	}
	setVersion(&u, 4)
	setVariant(&u, 1)

	ret := uuid2String(&u)

	if len(keepLine) == 0 || !keepLine[0] {
		ret = UuidRemoveLine(ret)
	}

	// ret := uuid.NewV4().String()

	// if len(keepLine) == 0 || !keepLine[0] {
	// 	ret = UuidRemoveLine(ret)
	// }

	return ret
}
