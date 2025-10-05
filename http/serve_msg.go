/*
Package http implements http controller for long polling.
*/
package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"log"

	"turbo_snail/broker"
)


type outMsg struct {
	Priority int8        `json:"priority"`
	Data   any	`json:"data"`
}

type ACKNACKSuccessMsg struct {
	Msg string `json:"msg"`
}



func SendMsg(w http.ResponseWriter, r *http.Request){
	turboSnailBroker := broker.Get() 
	vars := mux.Vars(r)
	trackName := vars["track"]
	msg := turboSnailBroker.GetMessage(trackName)
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


func ACKHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)	
	trackName := vars["track"]
	msgID := vars["msgId"]

	_track:= broker.Get().Tracks[trackName]
	msgUUID , err := uuid.FromBytes([]byte(msgID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err :=	_track.ACKMessage(msgUUID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	successMsg := &ACKNACKSuccessMsg{ Msg : fmt.Sprintf("ACK recieved for the msg id: %s", msgID)}

	if err := json.NewEncoder(w).Encode(successMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func NACKHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackName := vars["track"]
	msgID := vars["msgID"]

	_track := broker.Get().Tracks[trackName]

	msgIDBytes, err := uuid.FromBytes([]byte(msgID))
	if err != nil {
		// fmt.Errorf("invalid msg uuid is provided %s", msgID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := _track.NACKMessage(msgIDBytes); err != nil {
		
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	successMsg := &ACKNACKSuccessMsg{ Msg : fmt.Sprintf("NACK recieved for the msg id: %s", msgID)}
	if err := json.NewEncoder(w).Encode(successMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
}

func StartServer(httpPort string, wg *sync.WaitGroup){
	defer wg.Done()
	router := mux.NewRouter()
	router.HandleFunc("/{track}/message", SendMsg ).Methods("GET")
	router.HandleFunc("/{track}/message/{msgId}/ack", ACKHandler).Methods("POST")
	router.HandleFunc("/{track}/message/{msgId}/nack", NACKHandler).Methods("POST")

	log.Printf("Http server starting on port: %s\n", httpPort)
	if err := http.ListenAndServe(":"+httpPort, router); err != nil {
		log.Fatalf("error occured while starting http server \n %s \n", err.Error())
	}

} 
