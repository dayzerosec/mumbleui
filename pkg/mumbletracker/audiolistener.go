package mumbletracker

import (
	"layeh.com/gumble/gumble"
	"time"
)

type AudioListener struct {
	Frequency time.Duration
	OnStartSpeaking func(*gumble.User)
	OnStopSpeaking func(*gumble.User)
}

func (m *AudioListener) OnAudioStream(e *gumble.AudioStreamEvent) {
	go m.ProcessStream(e.User, e.C)
}

func (m *AudioListener) ProcessStream(user *gumble.User, ch <-chan *gumble.AudioPacket) {
	ticker := time.NewTicker(m.Frequency)

	speaking := false
	for {
		select {
		case <- ch:
			ticker.Reset(m.Frequency)
			if !speaking {
				speaking = true
				m.OnStartSpeaking(user)
			}
		case <- ticker.C:
			ticker.Stop()
			if speaking {
				speaking = false
				m.OnStopSpeaking(user)
			}
		}
	}
}
