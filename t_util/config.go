package t_util

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigFunc interface {
	ReadConfig(filename string) int
	GetService() string
}

type ConfigInfo struct {
	Service string `json:"service"`
}

func init() {
}

func (conf *ConfigInfo) ReadConfig(filename string) int {

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

func (conf ConfigInfo) GetService() string {
	return conf.Service
}