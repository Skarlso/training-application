package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type persister struct {
	config *appConfig
}

var dirPath = "./data/"
var metaInfoFilePath = dirPath + "metainfo.txt"

func newPersister(appConfig *appConfig) (*persister, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, err
	}
	return &persister{
		config: appConfig,
	}, nil
}

func (p *persister) writeMetaInfo() {

	metaInfoFile, err := os.OpenFile(metaInfoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("cannot open file %s: %v\n", metaInfoFilePath, err)
	}
	defer metaInfoFile.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// TODO write application startup and shutdown into metainfo file

	for {
		select {
		case <-ticker.C:
			podName := os.Getenv("POD_NAME")
			podIP := os.Getenv("POD_IP")
			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			_, err = fmt.Fprintf(metaInfoFile, "%s pod name %s, pod ip %s\n", timeStamp, podName, podIP)
			if err != nil {
				log.Errorf("cannot append to file %s: %v\n", metaInfoFilePath, err)
				return
			}
		}
	}
}
