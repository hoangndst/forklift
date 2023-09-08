package gcp

import (
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/storage"
	"context"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"io"
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

func (c *Client) connectStorageServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.storageService == nil {
		storageClient, err := storage.NewClient(c.ctx)
		if err != nil {
			c.storageService = storageClient
		}
	}
	return
}

func (c *Client) VMStart(project string, zone string, instance string) (err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return err
	}
	req := &computepb.StartInstanceRequest{
		Project:  project,
		Instance: instance,
		Zone:     zone,
	}
	_, err = c.computeService.Start(c.ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) VMStop(project string, zone string, instance string) (err error) {
	err = c.connectComputeServiceAPI()

	if err != nil {
		return err
	}
	req := &computepb.StopInstanceRequest{
		Project:  project,
		Instance: instance,
		Zone:     zone,
	}
	_, err = c.computeService.Stop(c.ctx, req)
	if err != nil {
		return err
	}
	return
}

func (c *Client) DownloadImageFromBucket(bucketName string, objectName string) (reader io.ReadCloser, err error) {
	err = c.connectStorageServiceAPI()
	if err != nil {
		return
	}
	bucket := c.storageService.Bucket(bucketName)
	object := bucket.Object(objectName)
	reader, err = object.NewReader(c.ctx)
	return
}
