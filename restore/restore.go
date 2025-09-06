package restore

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"turbo_snail/broker"
	"turbo_snail/config"
	"turbo_snail/log_entry"
)

func Init(turboSnailBroker *broker.Broker) error {
	walDir := config.WAL_DIR
	walFiles, err := os.ReadDir(walDir)
	if err != nil {
		log.Fatalf("Error occured while checking for persisted msgs during start up. \n %s", err.Error())
		return err
	}
	return buildTracks(turboSnailBroker, walFiles, walDir)
}

func buildTracks(turboSnailBroker *broker.Broker, walFiles []os.DirEntry, walDir string) error {
	var err error
	for _, walFile := range walFiles {

		fullFileName := fmt.Sprintf("%s/%s", walDir, walFile.Name())
		if !walFile.IsDir() && strings.Contains(fullFileName, ".log") {
			file, err := os.Open(fullFileName)
			if err != nil {
				log.Fatalf("Error occured while trying to read file %s \n %s", fullFileName, err.Error())
			}
			decoder := gob.NewDecoder(file)

			for {
				entry := log_entry.LogEntry{}
				err := decoder.Decode(&entry)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("Error occured while decoding logs from %s file.\n %s", fullFileName, err.Error())
				}

				if !entry.IsDeleted() {
					turboSnailBroker.AppendMsg(getActualTrackName(walFile.Name()), entry.Message)
				}
			}

			// read loggedBytes , parse each gob message and also ignore the deleted messages , remove them from in memory queue
			// 1st parse all the messages from the wal file

		}
	}

	return err
}

func getActualTrackName(fullFileName string) string {
	index := strings.LastIndex(fullFileName, ".log")
	return fullFileName[:index]
}
