package redis

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/pkrss/go-utils/profile"
	// _ "github.com/astaxie/beego/cache/redis"
)

var cc *RedisCache

// func init() {
// 	cc = nil

// 	initRedis()
// }

func InitRedis() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("initial redis error caught: %v\n", r)
			cc = nil
		}
	}()

	con := profile.ProfileReadString("MY_REDIS_IP", "127.0.0.1")
	if !strings.Contains(con, ":") {
		con += ":" + profile.ProfileReadString("MY_REDIS_PORT", "6379")
	}

	key := profile.ProfileReadString("MY_REDIS_PASSWORD")

	constr := `{"conn":"` + con + `"`

	if len(key) > 0 {
		constr += `,"password":"` + key + `"`
	}

	constr += `,"dbNum":0}`

	log.Println("redis connect:" + con)

	// cc, err = cache.NewCache("redis", constr)
	// if err != nil {
	// 	log.Println("redis error:" + err.Error())
	// }

	cc = RedisNewCache(con, key)
}

func SetCache(key string, value string, timeout int) error {

	if cc == nil {
		InitRedis()
	}

	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("set cache error caught: %v\n", r)
			cc = nil
		}
	}()
	timeouts := time.Duration(timeout) * time.Second
	err := cc.Put(key, value, timeouts)
	if err != nil {
		//log.Println("Cache失败，key:", key)
		return err
	} else {
		//log.Println("Cache成功，key:", key)
		return nil
	}
}

func GetCache(key string) string {
	if cc == nil {
		InitRedis()
	}

	if cc == nil {
		log.Println("get cache error cc is nil")
		return ""
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		return ""
	}

	switch v := data.(type) {
	case []uint8:
		return string(v)
	case string:
		return v
	default:
		log.Printf("get cache unknown type: %v\n", v)
		return ""
	}
	return data.(string)
}

func DelCache(key string) error {
	if cc == nil {
		InitRedis()
	}

	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	err := cc.Delete(key)
	if err != nil {
		return errors.New("Cache delete faield")
	} else {
		//fmt.Println("删除Cache成功 " + key)
		return nil
	}
}

// // --------------------
// // Encode
// // 用gob进行数据编码
// //
// func Encode(data interface{}) ([]byte, error) {
// 	buf := bytes.NewBuffer(nil)
// 	enc := gob.NewEncoder(buf)
// 	err := enc.Encode(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// // -------------------
// // Decode
// // 用gob进行数据解码
// //
// func Decode(data []byte, to interface{}) error {
// 	buf := bytes.NewBuffer(data)
// 	dec := gob.NewDecoder(buf)
// 	return dec.Decode(to)
// }
