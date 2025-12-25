import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func sendAsync(event Event) {
	go func() {
		data, _ := json.Marshal(event)
		apiURL := "https://termstat-api-production.up.railway.app/v1/ping"
		script := fmt.Sprintf(`curl -s -X POST %s -H "Content-Type: application/json" -d '%s' > /dev/null 2>&1 &`, apiURL, string(data))

		cmd := exec.Command("sh", "-c", script)
		_ = cmd.Start()
	}()
}