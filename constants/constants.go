package constants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Servers []ServerConfig `json:"servers"`
	Proxy   string         `json:"proxy"`
}

type ServerConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

type Backend struct {
	URL           string `json:"url"`
	Healthy       bool
	Connections   int
	ResponseTime  time.Duration
	Weight        int `json:"weight"`
	WeightedScore float64
}

func (b *Backend) IncrementConnections() {
	b.Connections++
	fmt.Println("Connections increment", b.Connections)
}
func (b *Backend) DecrementConnections() {
	fmt.Println("Connections decrement", b.Connections)
	b.Connections--
}

var (
	RoundRobbin         = "RoundRobbin"
	LeastConnections    = "LeastConnections"
	LeastTime           = "LeastTime"
	WeightedRoundRobbin = "WeightedRoundRobbin"
)

func ReadConfig(filename string) (*Config, error) {
	// Read JSON file
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var config Config

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &config)
	// Populate Backend struct array
	var backends []Backend
	for _, server := range config.Servers {
		backend := Backend{
			URL:     server.URL,
			Weight:  server.Weight,
			Healthy: true, // Initialize as healthy by default
		}
		backends = append(backends, backend)
	}

	return &config, nil
}
