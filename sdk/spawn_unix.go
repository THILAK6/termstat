// go:build !windows
// +build !windows

package sdk

import (
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
)

func spawnBackground(event Event) {
		data, _ := json.Marshal(event)
		exe, _ := os.Executable()

		cmd := exec.Command(exe, InternalFlag, encoded)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		_ = cmd.Start()
	} 		 