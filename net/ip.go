package net

import (
	"strconv"
	"strings"
)

// func Inet_ntoa(ipnr int64) net.IP {
//     var bytes [4]byte
//     bytes[0] = byte(ipnr & 0xFF)
//     bytes[1] = byte((ipnr >> 8) & 0xFF)
//     bytes[2] = byte((ipnr >> 16) & 0xFF)
//     bytes[3] = byte((ipnr >> 24) & 0xFF)

//     return net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0])
// }

// // Convert net.IP to int64 ,  http://www.sharejs.com
// func Inet_aton(ipnr net.IP) int {
// 	return Ip_s2i(ipnr.String())
// }

func Ip_i2s(i int) string {
	var bytes [4]byte
	bytes[0] = byte(i & 0xFF)
	bytes[1] = byte((i >> 8) & 0xFF)
	bytes[2] = byte((i >> 16) & 0xFF)
	bytes[3] = byte((i >> 24) & 0xFF)

	return strconv.Itoa(int(bytes[3])) + "." + strconv.Itoa(int(bytes[2])) + "." + strconv.Itoa(int(bytes[1])) + "." + strconv.Itoa(int(bytes[0]))
}

func Ip_s2i(s string) int {
	bits := strings.Split(s, ".")
	if len(bits) < 4 {
		return 0
	}

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int

	sum += int(b0) << 24
	sum += int(b1) << 16
	sum += int(b2) << 8
	sum += int(b3)

	return sum
}
