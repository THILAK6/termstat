package sdk

import (
	"os"
	"runtime"
	"time"
	"encoding/json"
	"fmt"
	"os/exec"
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

func sendAsync(event Event) {
		data, _ := json.Marshal(event)
		apiURL := "https://termstat-api-production.up.railway.app/v1/ping"
		script := fmt.Sprintf(`curl -s -X POST %s -H "Content-Type: application/json" -d '%s' > /dev/null 2>&1 &`, apiURL, string(data))

		cmd := exec.Command("sh", "-c", script)
		_ = cmd.Start()
}