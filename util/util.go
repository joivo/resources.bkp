package util

import log "github.com/sirupsen/logrus"

func HandleErr(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
