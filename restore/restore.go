package restore

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"strings"
	"turbo_snail/broker"
	"turbo_snail/log_entry"
)

func Init(turboSnailBroker *broker.Broker) error {
	walFiles, err := os.ReadDir("wal/")
	if err != nil {
		log.Fatalf("Error occured while checking for persisted msgs during start up. \n %s", err.Error())
		return err
	}
	return buildTracks(turboSnailBroker, walFiles)
}

func buildTracks(turboSnailBroker *broker.Broker, walFiles []os.DirEntry) error {
	var err error
	for _, walFile := range walFiles {
		fileName := walFile.Name()
		if !walFile.IsDir() && strings.Contains(fileName, ".log") {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			decoder := gob.NewDecoder(file)

			for {
				entry := log_entry.LogEntry{}
				err := decoder.Decode(&entry)
				log.Fatal(err)
				if err == io.EOF {
					break
				}

				if !entry.IsDeleted() {
					turboSnailBroker.AppendMsg(fileName, entry.Message)
				}
			}

			// read loggedBytes , parse each gob message and also ignore the deleted messages , remove them from in memory queue
			// 1st parse all the messages from the wal file

		}
	}

	return err
}
