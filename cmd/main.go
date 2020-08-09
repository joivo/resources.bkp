package main

import (
	"os"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/osbckp"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()

	log.Infoln("**** Starting Backup Service ****")

	osbckp.RegisterWorkers()
	osbckp.StartWorkers()

	log.Infoln("Exiting...")
	os.Exit(0)
}
