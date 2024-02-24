package healthcheck

import (
	// "fmt"
	// ""
	"net/http"

	"github.com/lokesh-katari/GoScale/constants"
)

func CheckHealth(backend *constants.Backend) bool {
	// Send a health check request to the backend server
	resp, err := http.Get(backend.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
