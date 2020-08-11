package main

import (
	"os"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/osbckp"
	"github.com/joivo/osbckp/util"

	"github.com/gophercloud/gophercloud"
	"github.com/nuveo/log"
)

func main() {
	log.Println("*** Starting OpenStack snapshots backup ***")

	config.LoadConfig()

	provider, err := osbckp.CreateClientProvider()
	util.HandleErr(err)

	regionName := config.GetOpenStackConfig().Clouds.OpenStack.RegionName
	endpointOpts := gophercloud.EndpointOpts{
		Region:       regionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	osbckp.RegisterWorker(osbckp.SnapshotWorkerCreator(provider, endpointOpts))
	osbckp.StartWorkers(config.FifteenDaysInMin, provider, endpointOpts)

	log.Println("Exiting...")
	os.Exit(0)
}
