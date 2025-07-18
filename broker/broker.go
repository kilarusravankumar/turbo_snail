package broker

import (
	"sync"

	"turbo_snail/message"
	"turbo_snail/track"
)

type Broker struct{
	Tracks map[string]*track.Track
	rwMutex sync.RWMutex
}

var turboSnailBroker *Broker
var once sync.Once

func CreateBroker() *Broker{
	once.Do(func(){
		turboSnailBroker = &Broker{
			Tracks : map[string]*track.Track{},
			rwMutex: sync.RWMutex{},
		}	
	})	
	return turboSnailBroker
	
}



func (b *Broker) AppendMsg(trackName string, data interface{}, priority int8){
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



func (b *Broker) GetMessage(trackName string) *message.Message {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()
	if raceTrack , exists := b.Tracks[trackName]; exists {
		return raceTrack.PopMessage()
	}
	return nil
}
