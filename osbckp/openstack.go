package osbckp

import (
	"strings"
	"sync"
	"time"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/nuveo/log"
)

func handleVolumeSnapshotResult(res snapshots.CreateResult, group *sync.WaitGroup, client *gophercloud.ServiceClient) {
	defer group.Done()
	snap, err := res.Extract()
	util.HandleErr(err)
	id := snap.ID
	log.Println("Handling snapshot result to volume with ID", id)
	log.Println("Snapshot initial status ", snap.Status)

	err = snapshots.WaitForStatus(client, id, config.UsefulVolumeStatus, 60)
	if err != nil {
		log.Println(err.Error())
		return
	}

	r := snapshots.Get(client, id)
	snap, err = r.Extract()

	util.HandleErr(err)
	log.Println("Snapshot status ", snap.Status)
}

func CreateVolumesSnapshots(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	log.Println("Creating volumes snapshots")

	bsV3, err := openstack.NewBlockStorageV3(provider, eopts)
	util.HandleErr(err)

	allPages, err := volumes.List(bsV3, volumes.ListOpts{
		Status:   config.UsefulVolumeStatus,
	}).AllPages()
	util.HandleErr(err)
	
	extractedVolumes, err := volumes.ExtractVolumes(allPages)
	util.HandleErr(err)
	
	log.Printf("%d volumes were found\n", len(extractedVolumes))

	wg := new(sync.WaitGroup)

	wg.Add(len(extractedVolumes))
	for _, v := range extractedVolumes {
		snapshotName := config.SnapshotSuffix + v.ID + "_" + time.Now().Format(config.DateLayout)
		desc := "Snapshot automatically created created by backup service"
		createSnapshotOpts := snapshots.CreateOpts{
			VolumeID:    v.ID,
			Force:       true,
			Name:        snapshotName,
			Description: desc,
		}
		log.Println("Snapshot name ", snapshotName)
		log.Printf("Sending request to snapshot for %s volume\n", v.Name)
		log.Printf("Creating snapshot of volume %s\n", v.ID)
		func(group *sync.WaitGroup) {
			r := snapshots.Create(bsV3, createSnapshotOpts)
			handleVolumeSnapshotResult(r, group, bsV3)
		}(wg)
	}

	wg.Wait()

	log.Println("Volumes snapshot finished")
}

func handleInstanceSnapshotResult(res servers.CreateImageResult, group *sync.WaitGroup, client *gophercloud.ServiceClient) {
	defer group.Done()
	var maxRetry  = 100

	id, err := res.ExtractImageID()
	util.HandleErr(err)

	for maxRetry != 0 {
		log.Println("Checking result of instance snapshot ", id)
		log.Println("Retry ", maxRetry)

		r := images.Get(client, id)

		var Response struct {
			Image struct{
				Status string `json:"status"`
			}
		}

		err := r.ExtractInto(&Response)
		util.HandleErr(err)

		currentStatus := strings.ToLower(Response.Image.Status)
		log.Printf("Image has status [%s]\n", currentStatus)

		if currentStatus == string(images.ImageStatusActive) {
			return
		}
		maxRetry = maxRetry - 1
		time.Sleep(config.PoolingInterval)
	}

	if maxRetry == 0 {
		log.Println("Worker exhausted, retry exceeded")
		return
	}
}

func CreateServersSnapshots(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	log.Println("Creating servers snapshots")

	computeV2, err := openstack.NewComputeV2(provider, eopts)
	util.HandleErr(err)

	allPages, err := servers.List(computeV2, servers.ListOpts{
		Status:     config.UsefulServerStatus,
		AllTenants: true,
	}).AllPages()
	util.HandleErr(err)

	extractedServers, err := servers.ExtractServers(allPages)
	util.HandleErr(err)

	log.Printf("%d instances were found\n", len(extractedServers))

	wg := new(sync.WaitGroup)

	wg.Add(len(extractedServers))
	for _, srv := range extractedServers {
		
		snapshotName := config.SnapshotSuffix + srv.Name + "_" + time.Now().Format(config.DateLayout)
		createImgOpts := servers.CreateImageOpts{
			Name: snapshotName,
		}
		log.Println("Snapshot name ", snapshotName)
		log.Printf("Sending request to build image for %s\n", srv.Name)
		log.Printf("Creating snapshot of server %s\n", srv.ID)
		srv := srv
		func(w *sync.WaitGroup) {
			group := new(sync.WaitGroup)
			group.Add(1)

			r := servers.CreateImage(computeV2, srv.ID, createImgOpts)
			handleInstanceSnapshotResult(r, group, computeV2)

			group.Wait()
			w.Done()
		}(wg)
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