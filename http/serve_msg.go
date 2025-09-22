package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"log"

	"turbo_snail/broker"
)


type outMsg struct {
	Priority int8        `json:"priority"`
	Data     interface{} `json:"data"`
}


func SendMsg(w http.ResponseWriter, r *http.Request){
	turboSnailBroker := broker.Get() 
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

func StartServer(HTTP_PORT string, wg *sync.WaitGroup){
	defer wg.Done()
	router := mux.NewRouter()
	router.HandleFunc("/{track}/message", SendMsg ).Methods("GET")

	log.Printf("Http server starting on port: %s\n", HTTP_PORT)
	if err := http.ListenAndServe(":"+HTTP_PORT, router); err != nil {
		log.Fatalf("error occured while starting http server \n %s \n", err.Error())
	}

} 
