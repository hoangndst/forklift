package gcp

import (
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/storage"
	"context"
	"os"
)

type Client struct {
	GoogleAuthPath  string
	ctx             context.Context
	computeMetadata *metadata.Client
	computeService  *compute.InstancesClient
	imageService    *compute.ImagesClient
	networkService  *compute.NetworksClient
	storageService  *storage.Client
}

func (c *Client) Authenticate() (err error) {
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", c.GoogleAuthPath)
	if err != nil {
		return
	}
	return
}

// Connect
func (c *Client) Connect() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	//userProjects := &[]Project{}
	//err = c.GetUserProjects(userProjects)
	return
}

func (c *Client) connectComputeMetadataAPI() (err error) {
	if c.computeMetadata == nil {

		if err != nil {
			return
		}
	}
	return
}

func (c *Client) connectComputeServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.computeService == nil {
		instancesClient, err := compute.NewInstancesRESTClient(c.ctx)
		if err != nil {
			c.computeService = instancesClient
		}
	}
	return
}
