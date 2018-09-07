package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	pkContainer "github.com/pkrss/go-utils/container"
)

// type ConfigStruct struct {
// 	AppName  string `json:"appname"`
// 	HttpPort int    `json:"httpport"`
// 	RunMode  string `json:"runmode"`
// 	Sys      struct {
// 		SmsCode struct {
// 			Hack string `json:"hack"`
// 		} `json:"smscode"`
// 	} `json:"sys"`
// }

var Config map[string]interface{}

var confFileRelPath = "./conf/"

// func init() {
// 	LoadConfigFile("./conf/config.json")
// }

func InitByConfigFile(confFile ...string) {
	f := ""
	if len(confFile) > 0 {
		f = confFile[0]

		p := strings.Index(f, "./conf/")
		if p >= 0 {
			confFileRelPath = f[0:p] + confFileRelPath
		}
	} else {
		f = "./conf/config.json"
	}
	LoadConfigFile(&Config, f)

	s, ok := getEnvString("config.json")
	if ok && (s != "") {
		var m2 map[string]interface{}
		e := json.Unmarshal([]byte(s), &m2)
		if e == nil {
			pkContainer.MapMerge(Config, m2)
		} else {
			log.Printf("merge from sys_env config.json error: %v\n", e)
		}
	}

	runmode := GetString("sys.runmode")
	if runmode == "" {
		return
	}
	f = strings.Replace(f, ".json", "-"+runmode+".json", -1)
	if _, e := os.Stat(f); e == nil {
		m2 := make(map[string]interface{}, 0)
		if nil == LoadConfigFile(&m2, f) {
			pkContainer.MapMerge(Config, m2)
		}
	}

}

func LoadConfigFile(ob interface{}, confFile string) error {
	content, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = json.Unmarshal(content, ob)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func getEnvString(key string) (ret string, ok bool) {
	v, ok := os.LookupEnv(key)

	if !ok && strings.Contains(key, ".") {
		key2 := strings.Replace(key, ".", "_", -1)
		v, ok = os.LookupEnv(key2)
	}

	return v, ok
}

func GetString(key string) string {
	if key == "" {
		return ""
	}
	v, ok := getEnvString(key)
	if ok {
		return v
	}

	if Config == nil {
		return ""
	}
	return pkContainer.MapGetStringValue(Config, key)
}

func GetArray(key string) (ret []interface{}) {
	if key == "" {
		return
	}
	v, ok := getEnvString(key)
	if ok {
		if v[0] == '[' {
			e := json.Unmarshal([]byte(v), &ret)
			if e == nil {
				return ret
			} else {
				log.Printf("GetArray %s %v\n", key, e)
			}
		}
		vv := strings.Split(v, ",")
		ret = make([]interface{}, len(vv))
		for i, j := range vv {
			ret[i] = j
		}
		return
	}

	if Config == nil {
		return
	}
	v2 := pkContainer.MapGetValue(Config, key)
	if v2 != nil {
		return v2.([]interface{})
	}
	return
}

func GetStringArray(key string) (ret []string) {
	v := GetArray(key)
	if v == nil {
		return
	}
	ret = make([]string, len(v))
	for i, j := range v {
		ret[i] = j.(string)
	}
	return
}

func GetSubObject(key string) (ret map[string]interface{}) {
	if key == "" {
		return
	}
	v, ok := getEnvString(key)
	if ok {

		e := json.Unmarshal([]byte(v), &ret)
		if e == nil {
			return ret
		} else {
			log.Printf("GetSubObject %s %v\n", key, e)
		}
	}

	if Config == nil {
		return
	}
	v2 := pkContainer.MapGetValue(Config, key)
	if v2 != nil {
		return v2.(map[string]interface{})
	}
	return
}

// ReadFileData fileName="exchanges/zb/zb_rate.html" => "./conf/exchanges/zb/zb_rate.html"
func ReadFileData(fileName string) ([]byte, error) {
	return ioutil.ReadFile(confFileRelPath + fileName)
}
