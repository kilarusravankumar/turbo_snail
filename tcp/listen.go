package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"turbo_snail/broker"
	"turbo_snail/message"
)

type incomingMessage struct {
	Priority  int8                   `json:"priority"`
	Data      map[string]interface{} `json:"data"`
	TrackName string                 `json:"track"`
}

func Listen(port string, turboSnailBroker *broker.Broker, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error occured while listening for TCP on port: %s\n", port)
		log.Fatal(err.Error())
	}

	defer listener.Close()
	defer wg.Done()

	log.Printf("server listening to %s\n", port)
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Fatalf("Error Occured While Accepting Connection: %s", err.Error())
			continue
		}

		go handleIncomingMessages(conn, turboSnailBroker)

		// and add a track to the broker if it is not present in the map
	}
}

func handleIncomingMessages(conn net.Conn, turboSnailBroker *broker.Broker) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Error occurred while reading from conn: %s", err.Error())
	}
	_msg := incomingMessage{}
	err = json.Unmarshal(buf[:n], &_msg)
	if err != nil {
		log.Fatalf("Error occured while parsing the json : %s\n", err.Error())
	}

	// now turboSnailBroker.AppendMsg(trackName,data, priority)

	// this is bad i need to create message.Message here and pass it to the log and turboSnailBroker.
	turboSnailMessage := message.New(_msg.Data, _msg.Priority)

	turboSnailBroker.AppendMsg(_msg.TrackName, turboSnailMessage)
	fmt.Printf("\n--> Track %s has %d messages\n", _msg.TrackName, turboSnailBroker.Tracks[_msg.TrackName].Len())
}
