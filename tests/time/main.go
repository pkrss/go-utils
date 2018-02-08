package main

import (
	"log"
	"strings"
	"time"
)

func parseTime(d string) *time.Time {
	fmt := "2006-01-02 15:04:05"
	if strings.ContainsAny(d, "+-") {
		fmt = "2006-01-02 15:04:05Z07"
	}
	t, e := time.Parse(fmt, d)
	if e != nil {
		log.Printf("JSONTime parse [%s] error:%s\n", d, e.Error())
	}
	return &t
}
func main() {
	s := "2018-02-08 11:04:11+08"

	t := parseTime(s)

	log.Printf("parse [%v]=%v\n", s, t)

}
