package main

import (
	"github.com/joivo/osbckp/config"
	"os"

	"github.com/joivo/osbckp/osbckp"
	"github.com/joivo/osbckp/util"

	"github.com/nuveo/log"
)

func createLogPath() (logFilePath string) {
	logFilePath = "logs"
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logFilePath, 0755)
		util.HandleErr(err)
	}
	return
}

func loadLogFile() (f *os.File, err error){
	path := createLogPath()
	f, err = os.OpenFile(path+"/osbckp.log", os.O_WRONLY|os.O_CREATE, 0755)

	return
}

func main() {
	log.Println("*** Starting OpenStack snapshots backup ***")

	config.LoadConfig()

	osbckp.RegisterWorker(osbckp.SnapshotWorker)
	osbckp.StartWorkers()

	log.Println("Exiting...")
	os.Exit(0)
}
