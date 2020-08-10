package main

import (
	"os"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/osbckp"
	"github.com/joivo/osbckp/util"

	"github.com/nuveo/log"
)

func loadLogFile() (f *os.File, err error) {
	const path = "/var/log/osbckp"
	util.CreatePathIfNotExist(path)
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
