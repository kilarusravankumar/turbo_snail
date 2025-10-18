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

	"github.com/google/uuid"
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

		if walFile.IsDir() {
			continue
		}
		trackName := getActualTrackName(walFile.Name())
		ackMsgMap := GetAckLogMap(walDir, trackName)
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
				if !ackMsgMap[entry.Message.ID] {
					turboSnailBroker.AppendMsg(trackName, entry.Message)
				}
			}

		}
	}

	return err
}

func GetAckLogMap(walDir string, trackName string) map[uuid.UUID]bool {
	ackLogMap := make(map[uuid.UUID]bool, 0)
	ackLogFileName := fmt.Sprintf("%s.ack.log", trackName)
	ackLogFullFileName := fmt.Sprintf("%s/%s", walDir, ackLogFileName)

	file, err := os.Open(ackLogFullFileName)
	if err != nil {
		log.Printf("Couldn't open the ACK log file %s for track: %s", ackLogFullFileName, trackName)
		log.Printf("err : %s", err.Error())
		return ackLogMap
	}

	decoder := gob.NewDecoder(file)
	for {
		entry := log_entry.ACKLogEntry{}
		err := decoder.Decode(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error occured while decoding the ACK Log file. \n %s", err.Error())
		}
		ackLogMap[entry.MsgID] = true
	}

	return ackLogMap
}

func getActualTrackName(fullFileName string) string {
	index := strings.LastIndex(fullFileName, ".log")
	log.Printf("fullFileName : %s", fullFileName)
	return fullFileName[:index]
}
