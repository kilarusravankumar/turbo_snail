package restore

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"turbo_snail/broker"
	"turbo_snail/message"

	"github.com/stretchr/testify/assert"
)

func TestBuildTracks(t *testing.T) {

	t.Run("reading from the wal files", func(t *testing.T) {
		testData := []message.Message{
			{
				Data:      map[string]interface{}{"task": "terminate the job"},
				Timestamp: time.Now().Unix(),
				Priority:  1,
			},
			{
				Data:      nil,
				Timestamp: time.Now().Unix(),
				Priority:  0,
			},
			{
				Data:      map[string]interface{}{"task": "fetch the weather."},
				Timestamp: time.Now().Unix(),
				Priority:  0,
			},
		}

		// create test.gob file and add sample data to the test.go
		filePaths := []string{
			"/home/kreten/code/turboSnail/wal/test.gob",
			"/home/kreten/code/turboSnail/wal/test1.gob",
			"/home/kreten/code/turboSnail/wal/test2.gob",
		}
		for _, filePath := range filePaths {
			// Open the file with O_CREATE and O_EXCL flags.
			// O_CREATE: Create the file if it does not exist.
			// O_EXCL: Return an error if O_CREATE is set and the file already exists.
			// 0644: File permissions (read/write for owner, read-only for others).
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err != nil {
				if os.IsExist(err) {
					log.Printf("File '%s' already exists.\n", filePath)
				} else {
					log.Fatalf("Error creating file: %v\n", err)
				}
				return
			}
			defer file.Close() // Ensure the file is closed when the function exits.

			log.Printf("File '%s' created successfully.\n", filePath)

			encoder := gob.NewEncoder(file)

			for _, tData := range testData {
				encoder.Encode(tData)
			}

		}

		// test restore functionality
		testBroker := broker.Get()

		walFiles, err := os.ReadDir("wal")
		if err != nil {
			log.Fatalf("Error occured while reading Directory \"wal\" : %s", err.Error())
		}

		buildTracks(testBroker, walFiles)

		allTrackNames := testBroker.GetAllTrackNames()

		assert.Equal(t, 3, len(allTrackNames), fmt.Sprintf("length of tracks should match. expected %d , got %d", 3, len(allTrackNames)))

	})
}
