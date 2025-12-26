package sdk

import (
	"runtime"
	"time"
	"bytes"
	"net/http"
	"os"
	"strings"
)

const InternalFlag = "--termstat-internal-ping"

type Event struct {
	APIKey		string            `json:"api_key"`
	Cmd	   	string            `json:"cmd"`
	Version	string            `json:"version"`
	ExitCode	int               `json:"exit_code"`
	OS 	  	string            `json:"os"`
	Arch	  	string            `json:"arch"`
	Duration	int64            `json:"duration"`
	Timestamp	string         `json:"timestamp"`
}

func Track(apiKey, version, cmd string) func(int){

	if os.Getenv("TERMSAT_DISABLE") == "1" || os.Getenv("TERMSAT_DISABLE") == "true" {
		return func(exitCode int){}
	}

	handleInternalPing()

	start := time.Now()
	return func(exitCode int){
		event := Event{
			APIKey:    apiKey,
			Cmd:       scrub(cmd),
			Version:   version,
			ExitCode:  exitCode,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Duration:  time.Since(start).Milliseconds(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		spawnBackground(event)
	}
}

func scrub(command string) string {
	home, _ := os.UserHomeDir()
	if home != "" {
		command =  strings.ReplaceAll(command, home, "~")
	}
	return command
}

func handleInternalPing(){
	args := os.Args
	for i, arg := range args {
		if arg == InternalFlag && i+1 < len(args) {
			payload := args[i+1]

			apiURL := "https://termstat-api-production.up.railway.app/v1/ping"
			client := &http.Client{Timeout: 5 * time.Second}
			_, _ = client.Post(apiURL, "application/json", bytes.NewBuffer([]byte(payload)))

			os.Exit(0)
		}
	}
}