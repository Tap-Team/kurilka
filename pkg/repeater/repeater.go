package repeater

import "time"

type Repeater interface {
	// function which returns true if need try again and false if operation not need retry
	Repeat(f func() bool)
}

func NewRepeater(n int, pause time.Duration) Repeater {
	if n < 1 {
		n = 1
	}
	return &repeater{times: n, pause: pause}
}

type repeater struct {
	times int
	pause time.Duration
}

func (r *repeater) Repeat(rf func() bool) {
	for i := 0; i < r.times; i++ {
		repeat := rf()
		if !repeat {
			break
		}
		time.Sleep(r.pause)
	}
}

type PauseSettings struct {
	i  int
	st map[int]time.Duration
}

func NewPauseSettings(pauses ...time.Duration) *PauseSettings {
	p := PauseSettings{st: make(map[int]time.Duration, len(pauses))}
	for _, pause := range pauses {
		p.Add(pause)
	}
	return &p
}

func (p *PauseSettings) Add(pause time.Duration) {
	p.st[p.i] = pause
	p.i++
}

func (p PauseSettings) Len() int {
	return p.i
}

type settingsRepeater struct {
	settings *PauseSettings
}

func (sr *settingsRepeater) Repeat(rf func() bool) {
	for i := 0; i < sr.settings.Len(); i++ {
		repeat := rf()
		if !repeat {
			break
		}
		time.Sleep(sr.settings.st[i])
	}
}

func NewSettingsRepeater(pauseSettings *PauseSettings) Repeater {
	return &settingsRepeater{settings: pauseSettings}
}
