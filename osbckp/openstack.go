package osbckp

import (
	"sync"
	"time"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	log "github.com/sirupsen/logrus"
)

func GetSnapshotResult() {

}

func CreateServersSnapshots(provider *gophercloud.ProviderClient) {
	log.Infoln("Creating snapshots")
	regionName := config.GetOpenStackConfig().Clouds.OpenStack.RegionName
	computeOpts := gophercloud.EndpointOpts{
		Region:       regionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	computeV2, err := openstack.NewComputeV2(provider, computeOpts)
	util.HandleErr(err)

	allPages, err := servers.List(computeV2, servers.ListOpts{
		AllTenants: true,
	}).AllPages()
	util.HandleErr(err)

	srvs, err := servers.ExtractServers(allPages)
	util.HandleErr(err)

	log.Infof("%d images were found", len(srvs))

	wg := new(sync.WaitGroup)

	wg.Add(len(srvs))
	for _, srv := range srvs {
		const layout = "2006-01-02"
		snapshotName := srv.Name + "_" + time.Now().Format(layout)
		createImgOpts := servers.CreateImageOpts{
			Name: snapshotName,
		}
		log.Infoln("Snapshot name", snapshotName)
		log.Infof("Sending request to build image for %s", srv.Name)
		log.Infof("Creating snapshot of server %s\n", srv.ID)
		go servers.CreateImage(computeV2, srv.ID, createImgOpts)
	}
	wg.Wait()

	log.Infof("Snapshot finished")
}

func CreateClientProvider() (*gophercloud.ProviderClient, error) {
	log.Info("Creating client provider")
	authURL := config.GetOpenStackConfig().Clouds.OpenStack.Auth.AuthUrl
	userID := config.GetOpenStackConfig().Clouds.OpenStack.Auth.UserID
	passwd := config.GetOpenStackConfig().Clouds.OpenStack.Auth.Password
	projID := config.GetOpenStackConfig().Clouds.OpenStack.Auth.ProjectID

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: authURL,
		UserID:           userID,
		Password:         passwd,
		Scope: &gophercloud.AuthScope{
			ProjectID: projID,
		},
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	return provider, err
}
