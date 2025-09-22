package broker

import (
	"sync"
	"turbo_snail/message"
	"turbo_snail/track"
)

type Broker struct {
	Tracks  map[string]*track.Track
	rwMutex sync.RWMutex
}

var (
	turboSnailBroker *Broker
	once             sync.Once
)

func Get() *Broker {
	once.Do(func() {
		turboSnailBroker = &Broker{
			Tracks:  map[string]*track.Track{},
			rwMutex: sync.RWMutex{},
		}
	})
	return turboSnailBroker
}

func (b *Broker) AppendMsg(trackName string, msg *message.Message) {
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
	raceTrack.AddMessage(msg)
}

func (b *Broker) GetMessage(trackName string) *message.Message {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()
	if raceTrack, exists := b.Tracks[trackName]; exists {
		/*
			before returning pop'd msg , store it in memory as pittedMsg and then after certain timeout
			reintroduce the message back to the priority Queue.
		*/
		return raceTrack.PopMessage()
	}
	return nil
}

func (b *Broker) GetAllTrackNames() []string {
	trackNames := []string{}

	for trackName := range b.Tracks {
		trackNames = append(trackNames, trackName)
	}

	return trackNames
}
