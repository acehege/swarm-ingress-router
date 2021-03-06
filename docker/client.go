package docker

import (
	"log"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"
)

// ServicePuller takes labels and returns matching docker services
type ServicePuller interface {
	GetServices(map[string]string) []swarm.Service
}

// Client holds the configration needed for a connection to Docker's API
type Client struct {
	socket     string
	apiVersion string
}

// GetServices returns all Docker services mathcing the labels giben
func (c Client) GetServices(filterList map[string]string) []swarm.Service {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to lookup services: ", r)
		}
	}()

	if err != nil {
		log.Print("Failed to lookup services: ", err)
		return []swarm.Service{}
	}

	filter := filters.NewArgs()
	for k, v := range filterList {
		filter.Add(k, v)
	}

	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{Filter: filter})
	if err != nil {
		log.Print("Failed to lookup services: ", err)
		return []swarm.Service{}
	}

	return services
}

// NewClient returns a new instance of the client
func NewClient() ServicePuller {
	return ServicePuller(Client{socket: "unix:///var/run/docker.sock", apiVersion: "v1.24"})
}
