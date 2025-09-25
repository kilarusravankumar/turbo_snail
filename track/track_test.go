package track

import (
	"testing"
	"time"
	"turbo_snail/message"

	"github.com/stretchr/testify/assert"
)

func TestInFlightMessagesGeneration(t *testing.T) {
	testTrack := New("TestingTrack", "/home/kreten/wal_logs")
	testDataSlice := []map[string]interface{}{
		{"dog": "bruno", "cat": "tommy", "mouse": "jerry"},
		{"duckName1": "ralph", "duckName2": "not Raplh", "duckName": "ralph's cousin"},
	}
	expectedInflightMsgs := make([]message.InFlightMessage, len(testDataSlice)) // testDataMsgs := message.N

	for i, data := range testDataSlice {
		msg := message.New(data, int8(i))
		expectedInflightMsgs[i] = message.InFlightMessage{Message: msg, ExpiresAt: time.Now().Add(time.Minute)}
		testTrack.AddMessage(msg)
	}

	for i := len(expectedInflightMsgs) - 1; i >= 0; i-- {

		msg := testTrack.PopMessage()
		assert.Equal(t, expectedInflightMsgs[i].Message, testTrack.InFlightMessages[msg.ID].Message, "both messages should match.")
	}

}
