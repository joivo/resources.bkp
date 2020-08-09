package main

import (
	"os"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/osbckp"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("**** Starting Backup Service ****")

	config.LoadConfig()

	osbckp.RegisterWorkers()
	osbckp.StartWorkers()

	log.Infoln("Exiting...")
	os.Exit(0)
}
