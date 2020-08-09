package osbckp

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/joivo/osbckp/config"
	"sync"
	"time"

	"github.com/joivo/osbckp/util"

	"github.com/nuveo/log"
	"github.com/robfig/cron/v3"
)

type worker func(wg *sync.WaitGroup)
type job func()

var (
	mu      = new(sync.Mutex)
	workers = make([]worker, 0)
)

func SnapshotJob() {
	log.Println("Starting SnapShot Job at", time.Now())

	provider, err := CreateClientProvider()
	util.HandleErr(err)

	regionName := config.GetOpenStackConfig().Clouds.OpenStack.RegionName
	computeOpts := gophercloud.EndpointOpts{
		Region:       regionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	CreateServersSnapshots(provider, computeOpts)
	CreateVolumesSnapshots(provider, computeOpts)
}

func SnapshotWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	c := cron.New()
	const fifteenDaysInMin = 360
	schedAt := fmt.Sprintf("@every %dh", fifteenDaysInMin)

	entryId, err := c.AddFunc(schedAt, SnapshotJob)
	util.HandleErr(err)

	c.Run()

	log.Printf("EntryID: %s \n", entryId)
}

func startHandle(fn job) {
	log.Println("Running first start job")
	fn()
}

func RegisterWorker(fn worker) {
	mu.Lock()
	log.Println("Registering worker")
	workers = append(workers, fn)
	mu.Unlock()
}

func StartWorkers() {
	startHandle(SnapshotJob)

	log.Println("Starting cron workers")

	wg := new(sync.WaitGroup)
	mu.Lock()
	wg.Add(len(workers))

	for _, w := range workers {
		go w(wg)
	}
	mu.Unlock()
	wg.Wait()
}
