package osbckp

import (
	"fmt"
	"sync"
	"time"

	"github.com/joivo/osbckp/config"
	"github.com/joivo/osbckp/util"

	"github.com/gophercloud/gophercloud"
	"github.com/nuveo/log"
	"github.com/robfig/cron/v3"
)

type worker func(wg *sync.WaitGroup)
type job func()

const successCode = 0

var workPoints = make(chan int)

var (
	mu      = new(sync.Mutex)
	workers = make([]worker, 0)
)

func SnapshotJob() {
	log.Printf("Starting Snapshot Job at [%s]\n", time.Now().Format(config.DateLayout))

	provider, err := CreateClientProvider()
	util.HandleErr(err)

	regionName := config.GetOpenStackConfig().Clouds.OpenStack.RegionName
	computeOpts := gophercloud.EndpointOpts{
		Region:       regionName,
		Availability: gophercloud.AvailabilityAdmin,
	}

	CreateVolumesSnapshots(provider, computeOpts)
	CreateServersSnapshots(provider, computeOpts)
}

func SnapshotWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	c := cron.New()

	schedAt := fmt.Sprintf("@every %dh", config.FifteenDaysInMin)

	entryId, err := c.AddFunc(schedAt, SnapshotJob)
	util.HandleErr(err)
	log.Printf("EntryID: [%s] \n", entryId)
	c.Run()
	log.Println("Snapshot Worker done")
	workPoints <- successCode
}

func jobHandle(fn job) {
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
	jobHandle(SnapshotJob)

	log.Printf("Workers waiting [%v] minutes to wake up again\n", config.FifteenDaysInMin)

	wg := new(sync.WaitGroup)
	mu.Lock()
	wg.Add(len(workers))

	for _, w := range workers {
		go w(wg)
	}
	mu.Unlock()
	wg.Wait()
}
