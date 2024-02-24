package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	// "encoding/json"

	"github.com/lokesh-katari/GoScale/algorithms"
	backend "github.com/lokesh-katari/GoScale/constants"
	"github.com/lokesh-katari/GoScale/healthcheck"
)

// Define a struct to represent a backend server
var (
	mutex           sync.Mutex
	backends        []*backend.Backend // Slice of backend servers
	nextServerIndex = 0
	CurrentWeight   = 0
)

var algos = backend.WeightedRoundRobbin

func main() {
	// Initialize backends
	data, err := backend.ReadConfig("config.json")

	if err != nil {
		fmt.Println("Error reading config file", err)
	}

	urlString := data.Proxy
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	_, port, err := net.SplitHostPort(parsedURL.Host)
	strPort := fmt.Sprintf(":%s", port)
	if err != nil {
		fmt.Println("Error parsing port:", err)
		return

	}
	backends = make([]*backend.Backend, len(data.Servers))
	for i, s := range data.Servers {
		backends[i] = &backend.Backend{
			URL:         s.URL,
			Weight:      s.Weight,
			Healthy:     true,
			Connections: 0,
		}
		fmt.Println("Backend server", backends[i])
	}
	// Start HTTP server
	http.HandleFunc("/", handler)
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "GoScale is running successfully OK")
	})
	fmt.Printf("Load balancer started on %s", port)
	log.Fatal(http.ListenAndServe(strPort, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	backend := selectHealthyBackend()
	if backend == nil {
		fmt.Println("No healthy backend servers available")
		http.Error(w, "No healthy backend servers available", http.StatusServiceUnavailable)
		return
	}
	fmt.Println("Selected backend server:", backend.URL)
	url, _ := url.Parse(backend.URL)

	// Increment the number of connections for the selected backend server

	backend.IncrementConnections()

	// Forward request to selected backend server
	//the below is used to implemnet the keep alive connection upto 100
	startTime := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	proxy.ServeHTTP(w, r)
	backend.ResponseTime = time.Since(startTime)
	fmt.Println("Response time", backend.ResponseTime)

	// Decrement the number of connections for the selected backend server
	defer backend.DecrementConnections()

}

func selectHealthyBackend() *backend.Backend {
	// Filter out unhealthy backends
	healthyBackends := make([]*backend.Backend, 0)
	for i := range backends {
		if backends[i].Healthy {
			healthyBackends = append(healthyBackends, backends[i])
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	switch algos {
	case backend.RoundRobbin:
		nextServerIndex = algorithms.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.LeastConnections:
		nextServerIndex = algorithms.LeastConnections(healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.WeightedRoundRobbin:
		// fmt.Println("CurrentWeight is the", CurrentWeight)
		nextServerIndex, CurrentWeight = algorithms.WeightedRoundRobbin(CurrentWeight, healthyBackends)
		return healthyBackends[nextServerIndex]
	case backend.LeastTime:
		nextServerIndex = algorithms.LeastTime(healthyBackends)
		return healthyBackends[nextServerIndex]

	default:
		// Use round-robin as the default algorithm
		nextServerIndex = algorithms.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]

	}

}

func healthCheck() {
	for {
		time.Sleep(10 * time.Second) // Check health every 10 seconds
		// Perform health check for each backend server
		for i := range backends {
			if !healthcheck.CheckHealth(backends[i]) {
				// Mark backend as unhealthy
				backends[i].Healthy = false
				fmt.Printf("Backend %s is unhealthy %d \n", backends[i].URL, backends[i].Connections)
			} else {
				// Mark backend as healthy
				backends[i].Healthy = true
				fmt.Printf("Backend %s is healthy %d \n", backends[i].URL, backends[i].Connections)
			}
		}
	}
}

func init() {
	// Start health check in the background
	go healthCheck()
}
