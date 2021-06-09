package main

import (
	"fmt"
	"tnh-simple-server/proc_EventHandler"
	"tnh-simple-server/proc_Test"
	"tnh-simple-server/t_util"
)

func main() {
	var mainconfig t_util.ConfigInfo

	mainconfig.ReadConfig("config.json")

	switch mainconfig.GetService() {
	case "test":
		proc_Test.Main()
	case "eventhandler":
		proc_EventHandler.Main()
	default:
		fmt.Println("Service Error. Exit.")
	}

}
