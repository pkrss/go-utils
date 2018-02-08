package redis

import (
	"errors"
	"log"
	"strconv"
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

type ConnectOption struct {
	Password string
	SrvHost  string
	SrvPort  string
	DbNum    int
}

func ConnectOptionFromProfile() *ConnectOption {
	ret := ConnectOption{}
	ret.SrvHost = profile.ProfileReadString("MY_REDIS_IP", "127.0.0.1")
	if strings.Contains(ret.SrvHost, ":") {
		ss := strings.Split(ret.SrvHost, ":")
		if len(ss) > 1 {
			ret.SrvHost = ss[0]
			ret.SrvPort = ss[1]
		}
	} else {
		ret.SrvPort = profile.ProfileReadString("MY_REDIS_PORT", "6379")
	}

	ret.Password = profile.ProfileReadString("MY_REDIS_PASSWORD")
	ret.DbNum = profile.ProfileReadInt("MY_REDIS_DBNUM")

	return &ret
}

func InitRedis(opts ...*ConnectOption) {

	var opt *ConnectOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt == nil {
		opt = ConnectOptionFromProfile()
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("initial redis error caught: %v\n", r)
			cc = nil
		}
	}()

	con := opt.SrvHost

	if opt.SrvPort != "" {
		con += ":" + opt.SrvPort
	}

	constr := `{"conn":"` + con + `"`

	if opt.Password != "" {
		constr += `,"password":"` + opt.Password + `"`
	}

	constr += `,"dbNum":` + strconv.Itoa(opt.DbNum) + `}`

	log.Println("redis connect:" + con)

	// cc, err = cache.NewCache("redis", constr)
	// if err != nil {
	// 	log.Println("redis error:" + err.Error())
	// }

	cc = RedisNewCache(con, opt.Password)
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
