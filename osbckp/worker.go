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
	"github.com/zloylos/grsync"
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

func BackupJob() {
	log.Printf("Starting Backup Job at [%s]\n", time.Now().Format(config.DateLayout))
	CreateBackup()
}

func SharePointBackupWorker(wg *sync.WaitGroup) {

}

func SyncBackupWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Sync Backup Worker Started")

	nwg := new(sync.WaitGroup)
	src := config.GetRSyncConfig().Source
	dest := config.GetRSyncConfig().Destination
	rsh := config.GetRSyncConfig().RSH

	task := grsync.NewTask(
		src,
		dest,
		grsync.RsyncOptions{
			Verbose:       true,
			Checksum:      true,
			Recursive:     true,
			Compress:      true,
			HumanReadable: true,
			Progress:      true,
			Rsh:           rsh,
		},
	)
	nwg.Add(1)

	go func(w *sync.WaitGroup) {
		for {
			state := task.State()
			log.Printf(
				"progress: %.2f / rem. %d / tot. %d / sp. %s \n",
				state.Progress,
				state.Remain,
				state.Total,
				state.Speed,
			)
			time.Sleep(time.Second)
			if state.Progress == float64(100) {
				break
			}
		}
		w.Done()
	}(nwg)

	if err := task.Run(); err != nil {
		util.HandleErr(err)
	}

	log.Println(task.Log())
	nwg.Wait()
	log.Println("Backup finished")
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
	// jobHandle(SnapshotJob)
	// jobHandle(BackupJob)

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
