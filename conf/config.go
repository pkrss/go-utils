package conf

import (
	"encoding/json"
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

// func init() {
// 	LoadConfigFile("./conf/config.json")
// }

func InitByConfigFile(confFile ...string) {
	f := ""
	if len(confFile) > 0 {
		f = confFile[0]
	} else {
		f = "./conf/config.json"
	}
	LoadConfigFile(&Config, f)
	runmode := GetString("runmode")
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
	configFile, err := os.Open(confFile)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(ob)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func getEnvString(key string) (ret string, ok bool) {
	return os.LookupEnv(key)
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
