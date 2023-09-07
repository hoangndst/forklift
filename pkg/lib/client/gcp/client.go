package gcp

import (
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/storage"
)

type Client struct {
	computeService compute.InstancesClient
	imageService   compute.ImagesClient
	networkService compute.NetworksClient
	storageService storage.Client
}

func (c *Client) Authenticate() (err error) {
	return
}
