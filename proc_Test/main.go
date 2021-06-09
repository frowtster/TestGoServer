package proc_Test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tnh-simple-server/router"
	"tnh-simple-server/t_util"
)

type ConfigTest struct {
	t_util.ConfigInfo
	ListenHost string `json:"listenHost"`
	ListenPort int    `json:"listenPort"`
}

var testconfig ConfigTest

func (conf *ConfigTest) ReadConfig(filename string) int {

	data, err := os.Open(filename)
	if err != nil {
		fmt.Println("Open", err)
		return -1
	}
	defer data.Close()
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Decode", err)
		return -1
	}

	fmt.Println(conf)

	return 1
}

func GetConnectStr() string {
	var ret string
	ret = fmt.Sprintf("%s:%d", testconfig.ListenHost, testconfig.ListenPort)
	return ret
}

func Main() {
	testconfig.ReadConfig("config.json")

	t_util.Log.Info("Starting server on " + GetConnectStr())
	http.ListenAndServe(GetConnectStr(), router.Router)

}
