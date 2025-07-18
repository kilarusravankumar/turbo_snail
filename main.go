package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "math"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const(
	PORT = "7777"
	HTTP_PORT = "7000"
)

func main() {
	// recieve msgs from tcp connection 
	
	turboSnailBroker := CreateBroker() 

	//start http server
	go startHttpServer(turboSnailBroker)

	// on each tcp message , look for the track 

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s",PORT))
	if err != nil {
		log.Fatalf("Error occured while listening for TCP on port: %s\n", PORT)
		log.Fatal(err.Error())
	}

	defer listener.Close()

	fmt.Printf("server listening to %s\n", PORT)
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
	Data interface{} `json:"data"`
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
	fmt.Printf("\n--> Track %s has %d messages\n",_msg.TrackName, turboSnailBroker.Tracks[_msg.TrackName].Len())
	
}

type outMsg struct {
	Priority int8 `json:"priority"`
	Data interface{} `json:"data"`
}

func startHttpServer(turboSnailBroker *Broker) {
	router := mux.NewRouter()

	serveMsg := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		trackName := vars["track"]
		msg := turboSnailBroker.GetMessage(trackName)
		fmt.Println(msg)
		// retry := 1.0 
		for msg == nil {
			time.Sleep(time.Second * 1)
			msg = turboSnailBroker.GetMessage(trackName)
		}
		
		outGoingMsg := outMsg{Priority: msg.Priority, Data: msg.Data}

		if err := json.NewEncoder(w).Encode(outGoingMsg); err != nil {
		// Handle encoding errors
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	router.HandleFunc("/{track}/message", serveMsg).Methods("GET")
	if err := http.ListenAndServe(":"+HTTP_PORT, router) ; err != nil {
		log.Fatalf("error occured while starting http server \n %s \n", err.Error())
	}

	fmt.Printf("Http server started on port: %s", HTTP_PORT)
	
}
