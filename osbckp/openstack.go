package osbckp

import (
	"sync"
	"time"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/nuveo/log"
)

const dateLayout = "2006-01-02"

var (
	instancesRecorded = make(chan servers.Server)
	volumesRecorded = make(chan volumes.Volume)
)

func GetSnapshotResult() {

}

func CreateVolumesSnapshots(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	log.Println("Creating volumes snapshots")

	bsV3, err := openstack.NewBlockStorageV3(provider, eopts)
	util.HandleErr(err)

	allPages, err := volumes.List(bsV3, volumes.ListOpts{
		Status:   "available",
	}).AllPages()
	util.HandleErr(err)
	
	extractedVolumes, err := volumes.ExtractVolumes(allPages)
	util.HandleErr(err)
	
	log.Printf("%d volumes were found\n", len(extractedVolumes))

	wg := new(sync.WaitGroup)

	wg.Add(len(extractedVolumes))
	for _, v := range extractedVolumes {
		snapshotName := v.Name + "_" + time.Now().Format(dateLayout)
		desc := "Snapshot automatically created created by backup service"
		createSnapshotOpts := snapshots.CreateOpts{
			VolumeID:    v.ID,
			Force:       true,
			Name:        snapshotName,
			Description: desc,
		}
		log.Println("Snapshot name", snapshotName)
		log.Printf("Sending request to snapshot for %s volume\n", v.Name)
		log.Printf("Creating snapshot of volume %s\n", v.ID)
		go func(group *sync.WaitGroup) {
			snapshots.Create(bsV3, createSnapshotOpts)
			wg.Done()
		}(wg)
		volumesRecorded <- v
	}

	wg.Wait()

	log.Println("Volumes snapshot finished")
}

func CreateServersSnapshots(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	log.Println("Creating servers snapshots")

	computeV2, err := openstack.NewComputeV2(provider, eopts)
	util.HandleErr(err)

	allPages, err := servers.List(computeV2, servers.ListOpts{
		Status:     "ACTIVE",
		AllTenants: true,
	}).AllPages()
	util.HandleErr(err)

	extractedServers, err := servers.ExtractServers(allPages)
	util.HandleErr(err)

	log.Printf("%d instances were found\n", len(extractedServers))

	wg := new(sync.WaitGroup)

	wg.Add(len(extractedServers))
	for _, srv := range extractedServers {
		
		snapshotName := srv.Name + "_" + time.Now().Format(dateLayout)
		createImgOpts := servers.CreateImageOpts{
			Name: snapshotName,
		}
		log.Println("Snapshot name", snapshotName)
		log.Printf("Sending request to build image for %s\n", srv.Name)
		log.Printf("Creating snapshot of server %s\n", srv.ID)
		srv := srv
		go func(group *sync.WaitGroup) {
			servers.CreateImage(computeV2, srv.ID, createImgOpts)
			wg.Done()
		}(wg)
		instancesRecorded <- srv
	}

	wg.Wait()

	log.Println("Instances snapshot finished")
}

func CreateClientProvider() (*gophercloud.ProviderClient, error) {
	log.Println("Creating client provider")
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
