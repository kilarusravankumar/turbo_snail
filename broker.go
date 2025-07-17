package main

import (
	"sync"

	"turbo_snail/broker/message"
	"turbo_snail/broker/track"
)

type Broker struct{
	Tracks map[string]*track.Track
	rwMutex sync.RWMutex
}

func (b Broker) AppendMsg(trackName string, data []byte, priority int8){
	var raceTrack *track.Track
	b.rwMutex.RLock()
	if _, exists := b.Tracks[trackName]; exists {
		raceTrack = b.Tracks[trackName]	
	}
	b.rwMutex.RUnlock()

	if raceTrack == nil {
		b.rwMutex.Lock()
		raceTrack = track.New(trackName)
		b.Tracks[trackName] = raceTrack
		b.rwMutex.Unlock()

	}
	newMessage := message.New(data , priority )
	raceTrack.AppendMessage(newMessage)
}

