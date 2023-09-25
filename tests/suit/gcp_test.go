package suit

import (
	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"testing"
)

func listBuckets(w io.Writer, projectID string) ([]string, error) {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	var buckets []string
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
		fmt.Fprintf(w, "Bucket: %v\n", battrs.Name)
	}
	return buckets, nil
}

func listAllInstances(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()
	// Use the `MaxResults` parameter to limit the number of results that the API returns per response page.
	req := &computepb.AggregatedListInstancesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(3),
	}

	it := instancesClient.AggregatedList(ctx, req)
	fmt.Fprintf(w, "Instances found:\n")
	// Despite using the `MaxResults` parameter, you don't need to handle the pagination
	// yourself. The returned iterator object handles pagination
	// automatically, returning separated pages as you iterate over the results.
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			fmt.Fprintf(w, "%s\n", pair.Key)
			for _, instance := range instances {
				id, err := strconv.ParseUint("5346769639904429930", 10, 64)
				if err != nil {
					return err
				}
				fmt.Println(instance.GetId() == id)
				fmt.Fprintf(w, "- %s %s\n", instance.GetName(), instance.GetNetworkInterfaces())
			}
		}
	}
	return nil
}

func getInstances(projectID, instanceName string) error {
	// projectID := "your_project_id"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Instance: instanceName,
		Zone:     "us-central1-c",
	}

	instance, err := instancesClient.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get instance: %w", err)
	}

	fmt.Printf("Instance: %+v\n", instance.GetNetworkInterfaces())

	return nil
}

func getInstancesByID(projectID, instanceID string) error {
	// projectID := "your_project_id"
	// instanceID := "your_instance_id"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Instance: instanceID,
		Zone:     "us-central1-c",
	}

	instance, err := instancesClient.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get instance: %w", err)
	}

	fmt.Printf("Instance: %+v\n", instance)

	return nil
}

func startInstance(w io.Writer, projectID, zone, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.StartInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}
	op, err := instancesClient.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to start instance: %w", err)
	}

	// Wait for the create operation to complete.
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Instance started\n")

	return nil
}

func listImagesFromBucket(w io.Writer, projectID, bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	it := client.Bucket(bucketName).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects: %w", bucketName, err)
		}
		fmt.Fprintf(w, "Object: %v\n", attrs.Name)
	}
	return nil
}

func checkObjectIsReady(bucketName, objectName string) (exists bool, err error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	// check object is ready
	bucket := client.Bucket(bucketName)
	attrs, err := bucket.Object(objectName).Attrs(ctx)
	if err != nil {
		return false, fmt.Errorf("Object(%q).Attrs: %w", objectName, err)
	}
	return attrs != nil, nil
}

func getImageFromBucket(w io.Writer, projectID, bucketName, objectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", objectName, err)
	}
	defer rc.Close()

	fmt.Fprintf(w, "Object %v downloaded.\n", objectName)
	return nil
}

func imageList(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	req := &computepb.ListImagesRequest{
		Project: projectID,
	}

	it := imagesClient.List(ctx, req)
	fmt.Fprintf(w, "Images found:\n")
	for {
		image, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to list images: %w", err)
		}
		fmt.Fprintf(w, "- %s\n", image.GetSourceDisk())
	}
	return nil
}

func deleteImage(projectID, imageName string) error {
	// projectID := "your_project_id"
	// imageName := "your_image_name"

	ctx := context.Background()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	req := &computepb.DeleteImageRequest{
		Project: projectID,
		Image:   imageName,
	}

	op, err := imagesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete image: %w", err)
	}

	// Wait for the delete operation to complete.
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}
	fmt.Printf("Image deleted\n")

	return nil
}

func getInstanceBootDisk(projectID, zone, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get instance: %w", err)
	}

	fmt.Printf("Instance: %+v\n", instance)
	disks := instance.GetDisks()
	// get boot disk
	for _, disk := range disks {
		if disk.GetBoot() {
			fmt.Printf("Boot disk: %+v\n", disk)
			break
		}
	}

	return nil
}

func createBucket(w io.Writer, projectID, bucketName string) error {
	// projectID := "my-project-id"
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Bucket %v created\n", bucketName)
	return nil
}

func createImageFromDisk(
	w io.Writer,
	projectID, zone, sourceDiskName, imageName string,
	storageLocations []string,
	forceCreate bool,
) error {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// sourceDiskName := "your_disk_name"
	// imageName := "my_image"
	// // If storageLocations empty, automatically selects the closest one to the source
	storageLocations = []string{}
	// // If forceCreate is set to `true`, proceeds even if the disk is attached to
	// // a running instance. This may compromise integrity of the image!
	// forceCreate = false

	ctx := context.Background()
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()
	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// Get the source disk
	source_req := &computepb.GetDiskRequest{
		Disk:    sourceDiskName,
		Project: projectID,
		Zone:    zone,
	}

	disk, err := disksClient.Get(ctx, source_req)
	if err != nil {
		return fmt.Errorf("unable to get source disk: %w", err)
	}

	// Create the image
	req := computepb.InsertImageRequest{
		ForceCreate: &forceCreate,
		ImageResource: &computepb.Image{
			Name:             &imageName,
			SourceDisk:       disk.SelfLink,
			StorageLocations: storageLocations,
		},
		Project: projectID,
	}

	op, err := imagesClient.Insert(ctx, &req)

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk image %s created\n", imageName)

	return nil
}

func exportImageToBucket(projectID, imageName, imageType, bucketName, objectName string) error {
	ctx := context.Background()
	client, err := cloudbuild.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	// Create a Cloud Build build object.
	build := &cloudbuildpb.Build{
		Steps: []*cloudbuildpb.BuildStep{
			{
				Name: "gcr.io/compute-image-tools/gce_vm_image_export:release",
				Args: []string{
					"--timeout=7000s",
					"--source_image=image-1",
					"--client_id=api",
					"--format=qcow2",
					"--destination_uri=gs://forklift-1694331580/image-1.qcow2",
				},
			},
		},
	}
	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: projectID,
		Build:     build,
	}
	// Submit the build request.
	op, err := client.CreateBuild(ctx, req)
	if err != nil {
		log.Fatalf("Error creating Cloud Build: %v", err)
		return err
	}

	// Wait for the build operation to complete.
	_, err = op.Wait(ctx)
	if err != nil {
		log.Fatalf("Error waiting for build operation to complete: %v", err)
		return err
	}
	return nil
}

func getImageByName(projectID, imageName string) error {
	ctx := context.Background()
	client, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.GetImageRequest{
		Project: projectID,
		Image:   imageName,
	}

	image, err := client.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get image: %w", err)
	}

	fmt.Printf("Image: %+v\n", image.GetDescription())

	return nil
}

func listNetworks(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	client, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewNetworksRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListNetworksRequest{
		Project: projectID,
	}

	it := client.List(ctx, req)
	fmt.Fprintf(w, "Networks found:\n")
	for {
		network, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to list networks: %w", err)
		}
		fmt.Fprintf(w, "- %s\n", network.GetIPv4Range())
	}
	return nil
}

func listSubnetworks(w io.Writer, projectID string) error {
	// projectID := "your_project_id"
	ctx := context.Background()
	client, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewSubnetworksRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListSubnetworksRequest{
		Project: projectID,
		Region:  "us-central1",
	}

	it := client.List(ctx, req)
	fmt.Fprintf(w, "Subnetworks found:\n")
	for {
		subnetwork, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to list subnetworks: %w", err)
		}
		fmt.Fprintf(w, "- %s\n", subnetwork.GetName())
	}
	return nil
}

func TestListSubnetworks(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = listSubnetworks(w, "confident-sweep-395418")
	if err != nil {
		panic(err)
	}
}

func TestListNetworks(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = listNetworks(w, "confident-sweep-395418")
	if err != nil {
		panic(err)
	}
}

func TestListBuckets(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	buckets, err := listBuckets(w, "confident-sweep-395418")
	if err != nil {
		panic(err)
	}
	fmt.Println(buckets)
}

func TestListAllInstances(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = listAllInstances(w, "confident-sweep-395418")
	if err != nil {
		panic(err)
	}
}

func TestGetInstances(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	err = getInstances("confident-sweep-395418", "instance-1")
	if err != nil {
		panic(err)
	}
}

func TestGetInstancesByID(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	err = getInstancesByID("confident-sweep-395418", "5346769639904429930")
	if err != nil {
		panic(err)
	}
}

func TestStartInstance(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = startInstance(w, "confident-sweep-395418", "us-central1-c", "instance-1")
	if err != nil {
		panic(err)
	}
}

func TestListImagesFromBucket(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = listImagesFromBucket(w, "confident-sweep-395418", "example-image-1")
	if err != nil {
		panic(err)
	}
}

func TestCreateImageFromDisk(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = createImageFromDisk(w, "confident-sweep-395418", "us-central1-c", "instance-1", "image-go-test", []string{}, true)
	if err != nil {
		panic(err)
	}
}

func TestExportImageToBucket(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = exportImageToBucket("confident-sweep-395418", "image-go-test", "qcow2", "example-image-1", "image-go-test.qcow2")
}

func TestImageList(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := io.Writer(os.Stdout)
	err = imageList(w, "confident-sweep-395418")
	if err != nil {
		panic(err)
	}
}

func TestGetImageByName(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = getImageByName("confident-sweep-395418", "image-go-1")
	fmt.Printf("Image %s", computepb.Image_PENDING.String())
	if err != nil {
		panic(err)
	}
}

func TestDeleteImage(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = deleteImage("confident-sweep-395418", "image-go-test")
	if err != nil {
		panic(err)
	}
}

func TestGetInstanceBootDisk(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = getInstanceBootDisk("confident-sweep-395418", "us-central1-c", "instance-1")
	if err != nil {
		panic(err)
	}
}

func TestCreateBucket(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := io.Writer(os.Stdout)
	// generate bucket name from "forklift" + timestamp
	bucketName := fmt.Sprintf("forklift-%d", time.Now().Unix())
	err = createBucket(w, "confident-sweep-395418", bucketName)
	if err != nil {
		panic(err)
	}
}

func TestCheckObjectIsReady(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	exists, err := checkObjectIsReady("example-image-1", "image-go-1.qcow2")
	if err != nil {
		panic(err)
	}
	fmt.Println(exists)
}

func TestStuff(t *testing.T) {
	ctx := context.Background()
	// set context value
	ctx = context.WithValue(ctx, "bucketName", "example-image-1")
	ctx = context.WithValue(ctx, "objectName", "image-go-1.qcow2")
	bucketName := ctx.Value("bucketName")
	fmt.Println(bucketName.(string))
	objectName := ctx.Value("objectName")
	fmt.Println(objectName)
}
