package gcp

import (
	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"context"
	"fmt"
	liberr "github.com/konveyor/forklift-controller/pkg/lib/error"
	"github.com/konveyor/forklift-controller/pkg/lib/logging"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
	"io"
	core "k8s.io/api/core/v1"
	"log"
	"os"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/storage"
)

type Client struct {
	URL               string
	Options           map[string]string
	GoogleAuthPath    string
	ProjectID         string
	Zone              string
	BucketName        string
	Log               logging.LevelLogger
	ctx               context.Context
	computeService    *compute.InstancesClient
	imageService      *compute.ImagesClient
	networkService    *compute.NetworksClient
	diskService       *compute.DisksClient
	storageService    *storage.Client
	cloudBuildService *cloudbuild.Client
}

func (c *Client) LoadOptionsFromSecret(secret *core.Secret) {
	c.Options = make(map[string]string)
	for key, value := range secret.Data {
		c.Options[key] = string(value)
	}
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

func (c *Client) connectImageServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.imageService == nil {
		imagesClient, err := compute.NewImagesRESTClient(c.ctx)
		if err != nil {
			c.imageService = imagesClient
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

func (c *Client) connectNetworkServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.networkService == nil {
		networksClient, err := compute.NewNetworksRESTClient(c.ctx)
		if err != nil {
			c.networkService = networksClient
		}
	}
	return
}

func (c *Client) connectDiskServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.diskService == nil {
		disksClient, err := compute.NewDisksRESTClient(c.ctx)
		if err != nil {
			c.diskService = disksClient
		}
	}
	return
}

func (c *Client) connectCloudBuildServiceAPI() (err error) {
	err = c.Authenticate()
	if err != nil {
		return
	}
	if c.cloudBuildService == nil {
		cloudBuildClient, err := cloudbuild.NewClient(c.ctx)
		if err != nil {
			c.cloudBuildService = cloudBuildClient
		}
	}
	return
}

// List a resource.
func (c *Client) List(object interface{}, opts interface{}) (err error) {
	switch object.(type) {
	case *[]computepb.Instance:
		err = c.computeServiceAPI(object, opts)
	case *[]computepb.Image:
		err = c.imageServiceAPI(object, opts)
	case *[]computepb.Network:
		err = c.networkServiceAPI(object, opts)
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		return
	}
	return
}

// Get a resource.
func (c *Client) Get(object interface{}, ID string) (err error) {
	switch object.(type) {
	case *computepb.Instance:
		err = c.computeServiceAPI(object, &GetOpts{ID: ID})
	case *computepb.Image:
		err = c.imageServiceAPI(object, &GetOpts{ID: ID})
	case *computepb.Network:
		err = c.networkServiceAPI(object, &GetOpts{ID: ID})
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		err = liberr.Wrap(err, "trying to get object", object, "ID", ID)
	}
	return
}

func (c *Client) computeServiceAPI(object interface{}, opts interface{}) (err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	switch object.(type) {
	case *[]computepb.Instance, *computepb.Instance:
		err = c.vmAPI(object, opts)
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		return
	}
	return
}

func (c *Client) vmAPI(object interface{}, opts interface{}) (err error) {
	switch object.(type) {
	case *[]computepb.Instance:
		object := object.(*[]computepb.Instance)
		switch opts.(type) {
		case *VMListOpts:
			//opts := opts.(*VMListOpts)
			err = c.vmList(object)
		}
		if err != nil {
			return
		}
	case *computepb.Instance:
		object := object.(*computepb.Instance)
		switch opts.(type) {
		case *GetOpts:
			opts := opts.(*GetOpts)
			req := &computepb.GetInstanceRequest{
				Project:  c.ProjectID,
				Instance: opts.ID,
				Zone:     c.Zone,
			}
			instance, err := c.computeService.Get(c.ctx, req)
			if err != nil {
				return err
			}
			*object = *instance
		default:
			err = c.unsupportedTypeError(opts)
		}
	default:
		err = c.unsupportedTypeError(object)
	}
	return
}

func (c *Client) vmList(object *[]computepb.Instance) (err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.AggregatedListInstancesRequest{
		Project:    c.ProjectID,
		MaxResults: proto.Uint32(3),
	}
	it := c.computeService.AggregatedList(c.ctx, req)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				*object = append(*object, *instance)
			}
		}
	}
	return
}

func (c *Client) imageServiceAPI(object interface{}, opts interface{}) (err error) {
	err = c.connectImageServiceAPI()
	if err != nil {
		return
	}
	switch object.(type) {
	case *computepb.Image, *[]computepb.Image:
		err = c.imageAPI(object, opts)
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

func (c *Client) imageAPI(object interface{}, opts interface{}) (err error) {
	switch object.(type) {
	case *[]computepb.Image:
		object := object.(*[]computepb.Image)
		switch opts.(type) {
		case *ImageListOpts:
			err = c.imageList(object)
		default:
			err = c.unsupportedTypeError(opts)
		}
	case *computepb.Image:
		object := object.(*computepb.Image)
		switch opts.(type) {
		case *GetOpts:
			opts := opts.(*GetOpts)
			req := &computepb.GetImageRequest{
				Project: c.ProjectID,
				Image:   opts.ID,
			}
			image, err := c.imageService.Get(c.ctx, req)
			if err != nil {
				return
			}
			*object = *image
		case *DeleteOpts:
			req := &computepb.DeleteImageRequest{
				Project: c.ProjectID,
				Image:   object.GetName(),
			}
			_, err = c.imageService.Delete(c.ctx, req)
			if err != nil {
				return
			}
		default:
			err = c.unsupportedTypeError(opts)
		}
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		return
	}
	return
}

func (c *Client) imageList(object *[]computepb.Image) (err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.ListImagesRequest{
		Project: c.ProjectID,
	}
	it := c.imageService.List(c.ctx, req)
	for {
		image, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		*object = append(*object, *image)
	}
	return
}

func (c *Client) networkServiceAPI(object interface{}, opts interface{}) (err error) {
	err = c.connectNetworkServiceAPI()
	if err != nil {
		return
	}
	switch object.(type) {
	case *computepb.Network, *[]computepb.Network:
		err = c.networkAPI(object, opts)
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

func (c *Client) networkAPI(object interface{}, opts interface{}) (err error) {
	switch object.(type) {
	case *[]computepb.Network:
		object := object.(*[]computepb.Network)
		switch opts.(type) {
		case *NetworkListOpts:
			err = c.networkList(object)
		default:
			err = c.unsupportedTypeError(opts)
		}
	case *computepb.Network:
		object := object.(*computepb.Network)
		switch opts.(type) {
		case *GetOpts:
			opts := opts.(*GetOpts)
			req := &computepb.GetNetworkRequest{
				Project: c.ProjectID,
				Network: opts.ID,
			}
			network, err := c.networkService.Get(c.ctx, req)
			if err != nil {
				return
			}
			*object = *network
		default:
			err = c.unsupportedTypeError(opts)
		}
	default:
		err = c.unsupportedTypeError(object)
	}
	if err != nil {
		return
	}
	return
}

func (c *Client) networkList(object *[]computepb.Network) (err error) {
	err = c.connectNetworkServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.ListNetworksRequest{
		Project: c.ProjectID,
	}
	it := c.networkService.List(c.ctx, req)
	for {
		network, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		*object = append(*object, *network)
	}
	return
}

//func (c *Client) subnetList(object *[]computepb.Subnetwork) (err error) {
//	err = c.connectNetworkServiceAPI()
//	if err != nil {
//		return
//	}
//	req := &computepb.ListSubnetworksRequest{
//		Project: c.ProjectID,
//	}
//	it := c.networkService.List(c.ctx, req)
//	for {
//		subnetwork, err := it.Next()
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			return err
//		}
//		*object = append(*object, *subnetwork)
//	}
//	return
//}

func (c *Client) createBucket(bucketName string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	err := c.connectStorageServiceAPI()
	if err != nil {
		return err
	}
	bucket := c.storageService.Bucket(bucketName)
	if err := bucket.Create(c.ctx, c.ProjectID, nil); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	c.BucketName = bucketName
	return nil
}

func (c *Client) checkObjectExists(bucketName string, objectName string) (bool, error) {
	err := c.connectStorageServiceAPI()
	if err != nil {
		return false, err
	}
	bucket := c.storageService.Bucket(bucketName)
	object := bucket.Object(objectName)
	_, err = object.Attrs(c.ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) VMExportImageToBucket(imageName, imageType, objectName string) error {
	err := c.connectCloudBuildServiceAPI()
	if err != nil {
		return err
	}
	if c.BucketName == "" {
		c.BucketName = fmt.Sprintf("forklift-%s", time.Now().Format("20060102150405"))
		err = c.createBucket(c.BucketName)
		if err != nil {
			return err
		}
	} else {
		exists, err := c.checkObjectExists(c.BucketName, objectName)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
	}

	// Create a CloudBuild build object.
	destinationURI := fmt.Sprintf("gs://%s/%s.%s", c.BucketName, objectName, imageType)
	build := &cloudbuildpb.Build{
		Steps: []*cloudbuildpb.BuildStep{
			{
				Name: "gcr.io/compute-image-tools/gce_vm_image_export:release",
				Args: []string{
					"--timeout=7000s",
					"--source_image=" + imageName,
					"--client_id=api",
					"--format=qcow2",
					"--destination_uri=" + destinationURI,
				},
			},
		},
	}
	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: c.ProjectID,
		Build:     build,
	}
	// Submit the build request.
	op, err := c.cloudBuildService.CreateBuild(c.ctx, req)
	if err != nil {
		log.Fatalf("Error creating Cloud Build: %v", err)
		return err
	}

	// Wait for the build operation to complete.
	_, err = op.Wait(c.ctx)
	if err != nil {
		log.Fatalf("Error waiting for build operation to complete: %v", err)
		return err
	}
	return nil
}

func (c *Client) VMGetInstanceBootDisk(instanceName string) (disk *computepb.AttachedDisk, err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.GetInstanceRequest{
		Project:  c.ProjectID,
		Zone:     c.Zone,
		Instance: instanceName,
	}
	instance, err := c.computeService.Get(c.ctx, req)
	if err != nil {
		return
	}
	for _, disk := range instance.Disks {
		if disk.GetBoot() {
			return disk, nil
		}
	}
	return
}

func (c *Client) getInstanceBootDisk(instanceName string) (disk *computepb.AttachedDisk, err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.GetInstanceRequest{
		Project:  c.ProjectID,
		Zone:     c.Zone,
		Instance: instanceName,
	}
	instance, err := c.computeService.Get(c.ctx, req)
	if err != nil {
		return
	}
	for _, disk := range instance.Disks {
		if disk.GetBoot() {
			return disk, nil
		}
	}
	return
}

func (c *Client) VMCreateImageFromDisk(instanceName string, forceCreate bool) (imageName string, err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	err = c.connectDiskServiceAPI()
	if err != nil {
		return
	}
	// // If storageLocations empty, automatically selects the closest one to the source
	storageLocations := []string{}
	// // If forceCreate is set to `true`, proceeds even if the disk is attached to
	// // a running instance. This may compromise integrity of the image!
	// forceCreate = false`
	disk, err := c.getInstanceBootDisk(instanceName)
	if err != nil {
		return
	}
	sourceReq := &computepb.GetDiskRequest{
		Project: c.ProjectID,
		Disk:    disk.GetDeviceName(),
		Zone:    c.Zone,
	}
	imageName = fmt.Sprintf("%s-%s", instanceName, "forklift")
	sourceDisk, err := c.diskService.Get(c.ctx, sourceReq)
	if err != nil {
		return
	}
	req := &computepb.InsertImageRequest{
		Project: c.ProjectID,
		ImageResource: &computepb.Image{
			Name:             &imageName,
			SourceDisk:       sourceDisk.SelfLink,
			StorageLocations: storageLocations,
		},
		ForceCreate: &forceCreate,
	}
	op, err := c.imageService.Insert(c.ctx, req)

	if err = op.Wait(c.ctx); err != nil {
		return imageName, fmt.Errorf("unable to wait for the operation: %w", err)
	}
	return
}

func (c *Client) VMStatus(instance string) (status string, err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return
	}
	req := &computepb.GetInstanceRequest{
		Project:  c.ProjectID,
		Instance: instance,
		Zone:     c.Zone,
	}
	instanceInfo, err := c.computeService.Get(c.ctx, req)
	if err != nil {
		return
	}
	status = instanceInfo.GetStatus()
	return
}

func (c *Client) VMStart(instance string) (err error) {
	err = c.connectComputeServiceAPI()
	if err != nil {
		return err
	}
	req := &computepb.StartInstanceRequest{
		Project:  c.ProjectID,
		Instance: instance,
		Zone:     c.Zone,
	}
	op, err := c.computeService.Start(c.ctx, req)
	if err != nil {
		return err
	}
	// Wait for the operation to complete.
	if err = op.Wait(c.ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}
	return nil
}

func (c *Client) VMStop(instance string) (err error) {
	err = c.connectComputeServiceAPI()

	if err != nil {
		return err
	}
	req := &computepb.StopInstanceRequest{
		Project:  c.ProjectID,
		Instance: instance,
		Zone:     c.Zone,
	}
	op, err := c.computeService.Stop(c.ctx, req)
	if err != nil {
		return err
	}
	// Wait for the operation to complete.
	if err = op.Wait(c.ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
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

func (c *Client) unsupportedTypeError(object interface{}) (err error) {
	err = liberr.New(fmt.Sprintf("unsupported type %T", object))
	return
}

// buon qua di
