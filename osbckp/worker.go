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

type Worker func(wg *sync.WaitGroup)
type Job func()

var (
	mu      = new(sync.Mutex)
	workers = make([]Worker, 0)
)

func SnapshotJobCreator(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) Job {
	return func() {
		log.Printf("Starting Job to snapshot instances and volumes at [%s]\n", time.Now().Format(config.DateLayout))
		CreateVolumesSnapshots(provider, eopts)
		CreateServersSnapshots(provider, eopts)
		checkOldSnapshotsJobCreator(provider, eopts)
	}
}

func checkOldSnapshotsJobCreator(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	log.Println("Checking olds snapshots")
	DeleteOldSnapshots(provider, eopts)
}

func SnapshotWorkerCreator(provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) Worker {
	return func(wg *sync.WaitGroup) {
		defer wg.Done()
		c := cron.New()

		schedAt := fmt.Sprintf("@every %dh", config.WeekInHours)

		_, err := c.AddFunc(schedAt, SnapshotJobCreator(provider, eopts))
		util.HandleErr(err)
		c.Run()
		log.Println("Snapshot Worker done")
	}
}

func jobHandle(fn Job) {
	log.Println("Running first SnapshotJobCreator Job")
	fn()
}

func RegisterWorker(fn Worker) {
	mu.Lock()
	log.Println("Registering Worker")
	workers = append(workers, fn)
	mu.Unlock()
}

func StartWorkers(sleepTime int, provider *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) {
	jobHandle(SnapshotJobCreator(provider, eopts))

	log.Printf("Workers waiting [%v] minutes to wake up again\n", sleepTime)

	wg := new(sync.WaitGroup)
	mu.Lock()
	wg.Add(len(workers))

	for _, w := range workers {
		go w(wg)
	}
	mu.Unlock()
	wg.Wait()
}
