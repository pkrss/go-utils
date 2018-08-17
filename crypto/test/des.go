package main

import (
	"log"

	"github.com/pkrss/go-utils/crypto"
)

func main() {
	key := "pkrss-1"
	origData := "23"
	encrypted, e := crypto.DesEncrypt(origData, key)
	if e != nil {
		log.Println(e.Error())
		return
	}

	log.Println("target: " + encrypted)

	key2, e := crypto.DesDecrypt(encrypted, key)
	if e != nil {
		log.Println(e.Error())
		return
	}

	log.Println("orig:" + key2)

}
