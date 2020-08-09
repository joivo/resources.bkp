package osbckp

import (
	"fmt"
	"sync"
	"time"

	"github.com/joivo/osbckp/util"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type worker func(wg *sync.WaitGroup)
type job func()

var workers []worker

func SnapshotJob() {
	log.Infoln("Starting SnapShot Job", time.Now())

	provider, err := CreateClientProvider()
	util.HandleErr(err)

	CreateServersSnapshots(provider)
}

func SnapshotWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	c := cron.New()
	const fifteenDaysInMin = 360
	schedAt := fmt.Sprintf("@every %dh", fifteenDaysInMin)

	entryId, err := c.AddFunc(schedAt, SnapshotJob)
	util.HandleErr(err)

	c.Run()

	log.Infof("EntryID: %s \n", entryId)
}

func RunStartUpJob(fn job) {
	log.Infoln("Running first start job")
	fn()
}

func RegisterWorker(fn worker) {
	workers = append(workers, fn)
}

func RegisterWorkers() {
	log.Infoln("Registering workers...")

	RegisterWorker(SnapshotWorker)

	log.Infof("%d workers registered", len(workers))
}

func StartWorkers() {
	RunStartUpJob(SnapshotJob)

	log.Infoln("Starting workers")

	wg := new(sync.WaitGroup)
	wg.Add(len(workers))

	for _, w := range workers {
		go w(wg)
	}

	wg.Wait()
}
