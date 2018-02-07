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

func GetString(key string) string {
	if key == "" {
		return ""
	}
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	if Config == nil {
		return ""
	}
	return pkContainer.MapGetStringValue(Config, key)
}
