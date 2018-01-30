package redis

import (
	"fmt"
	"time"
)

func LockRun(lockerKey string, cb func(), timeoutSeconds int) error {

	for ; timeoutSeconds > 0; timeoutSeconds-- {

		var err error
		var v interface{}

		if v, err = cc.do("SET", lockerKey, "EX", timeoutSeconds, "NX"); err != nil {
			return err
		}

		var s string
		switch v2 := v.(type) {
		case []uint8:
			s = string(v2)
		case string:
			s = v2
		default:
			return fmt.Errorf("get cache unknown type: %v", v2)
		}

		if s == "OK" {
			cb()
			return nil
		}

		time.Sleep(time.Second * 1)
	}

	cb()

	cc.do("DELETE", lockerKey)

	return nil
}
