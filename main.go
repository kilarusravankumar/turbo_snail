package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

)

const(
	PORT = "7777"
)

func main() {
	// recieve msgs from tcp connection 
	turboSnailBroker := &Broker{}


	// on each tcp message , look for the track 

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s",PORT))
	if err != nil {
		log.Fatalf("Error occured while listening for TCP on port: %s\n", PORT)
		log.Fatal(err.Error())
	}

	defer listener.Close()

	fmt.Printf("server listening to %s", PORT)
	var conn net.Conn
	for{
		conn, err = listener.Accept()	
		if err != nil {
			log.Fatalf("Error Occured While Accepting Connection: %s", err.Error())
			continue
		}

		go handleConnection(conn, turboSnailBroker)

	// and add a track to the broker if it is not present in the map
	}

}

type incomingMessage struct {
	Priority int8 `json:"priority"`
	Data []byte  `json:"data"`
	TrackName string `json:"track"`
}

func handleConnection(conn net.Conn, turboSnailBroker *Broker) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err:= conn.Read(buf)
	if err != nil {
		log.Fatalf("Error occurred while reading from conn: %s", err.Error())
	}
	_msg := incomingMessage{}
	err = json.Unmarshal(buf[:n], &_msg)
	if err != nil {
		log.Fatalf("Error occured while parsing the json : %s\n", err.Error())
	}

	// now turboSnailBroker.AppendMsg(trackName,data, priority)
	turboSnailBroker.AppendMsg(_msg.TrackName, _msg.Data, _msg.Priority)
	fmt.Printf("Track %s has %d messages",_msg.TrackName, turboSnailBroker.Tracks[_msg.TrackName].Len())
	
}
