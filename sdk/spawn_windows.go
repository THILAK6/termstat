// go:build windows
package sdk

import(
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
)

func spawnBackground(event Event) {
		data, _ := json.Marshal(event)
		exe, _ := os.Executable()

		cmd := exec.Command(exe, InternalFlag, string(data))

		const (
			CREATE_NEW_PROCESS_GROUP = 0x00000200
			DETACHED_PROCESS         = 0x00000008
		)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS,
		}

		_ = cmd.Start()
	}