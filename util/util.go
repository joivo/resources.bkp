package util

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func HandleFatal(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func HandleErr(err error) {
	if err != nil {
		log.Errorf("error: %v", err)
	}
}

func CreatePathIfNotExist(pathDir string) {
	_, err := os.Stat(pathDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pathDir, 0755)
		HandleErr(err)
	}
	return
}
