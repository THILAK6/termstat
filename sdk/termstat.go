package sdk

import (
	"runtime"
	"time"
	"bytes"
	"net/http"
	"os"
	"strings"
	"encoding/base64"
)

const InternalFlag = "--termstat-internal-ping"

type Event struct {
	APIKey		string            `json:"api_key"`
	Cmd	   	string            `json:"cmd"`
	Version	string            `json:"version"`
	ExitCode	int               `json:"exit_code"`
	OS 	  	string            `json:"os"`
	Arch	  	string            `json:"arch"`
	DurationMs	int64            `json:"duration_ms"`
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
			DurationMs:  time.Since(start).Milliseconds(),
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
			encodedPayload := args[i+1]

			decoded, err := base64.StdEncoding.DecodeString(encodedPayload)
			if err != nil {
				os.Exit(1)
			}

			apiURL := "https://termstat-api-production.up.railway.app/v1/ping"
			client := &http.Client{Timeout: 5 * time.Second}
			_, _ = client.Post(apiURL, "application/json", bytes.NewBuffer([]byte(decoded)))
			os.Exit(0)
		}
	}
}