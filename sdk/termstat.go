package sdk

import (
	"os"
	"runtime"
	"time"
)

type Event struct {
	APIKey		string            `json:"api_key"`
	Cmd	   	string            `json:"cmd"`
	Version	string            `json:"version"`
	ExitCode	int               `json:"exit_code"`
	OS 	  	string            `json:"os"`
	Arch	  	string            `json:"arch"`
	Duration	time.Duration     `json:"duration"`
	Timestamp	time.Time         `json:"timestamp"`
}

func Track(apiKey, version, cmd string) func(init){
	start := time.Now()
	return func(exitCode int){
		event := Event{
			APIKey:    apiKey,
			Cmd:       cmd,
			Version:   version,
			ExitCode:  exitCode,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Duration:  time.Since(start).Milliseconds(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		sendAsync(event)
	}
}