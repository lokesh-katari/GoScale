# GoScale

### GoScale is a versatile load balancing application written in Go that provides various algorithms for distributing incoming traffic among backend servers. It supports the following algorithms:

- ### Round Robin (RR)
- ### Weighted Round Robin (WRR)
- ### Least Connections
- ### Least Time

# Features

- #### Dynamic Routing: GoScale dynamically routes incoming requests to healthy backend servers based on different load balancing algorithms.
- #### Health Checking: It regularly checks the health of backend servers and excludes unhealthy servers from the routing pool.
- #### Algorithm Options: Users can choose from multiple load balancing algorithms to optimize traffic distribution based on their specific requirements.
- #### Configuration: GoScale can be easily configured via a JSON file, allowing users to define backend servers, their weights, and other settings.
- #### Docker support for easy deployment and scalability.
- #### HTTP Keep-Alive Support: GoScale implements HTTP Keep-Alive to reuse backend connections until the timeout expires, improving performance and reducing latency.

## Usage

### Installation

To install GoScale , you need to have Go installed. Then, you can clone the repository and build the application using the following commands:

```bash
git clone https://github.com/lokesh-katari/GoScale
cd loadpulse
go build -o goscale cmd/main.go
```

## or

### To use GoScale, simply pull the Docker image from Docker Hub:

```bash
docker pull lokeshkatari/goscale

```

### Then, create a config.json file to configure the backend servers, their weights, and the load balancing algorithm. Finally, run the Docker container:

```bash
docker run -v /path/to/config.json:/app/config.json -p 8082:8082 lokeshkatari/goscale:latest

```

# Configuration

GoScale is configured via a config.json file. You can specify the backend servers, their weights, the load balancing algorithm, and other settings in this file.

```json
{
  "servers": [
    {
      "url": "http://localhost:8081",
      "weight": 2
    },
    {
      "url": "http://localhost:8082",
      "weight": 1
    }
  ],
  "algorithm": "WeightedRoundRobin",
  "proxy": "http://localhost:8080"
}
```

# Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

# License

This project is licensed under the MIT License - see the LICENSE file for details.
