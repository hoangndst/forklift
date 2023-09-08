package suit

import (
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
	"io"
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
				fmt.Fprintf(w, "- %s %d\n", instance.GetName(), instance.GetId())
			}
		}
	}
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

	_, err = instancesClient.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to start instance: %w", err)
	}

	// Wait for the create operation to complete.

	fmt.Fprintf(w, "Instance started\n")

	return nil
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

func TestStartInstance(t *testing.T) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hoangndst/Desktop/kubevirt/forklift/tests/suit/vdt.json")
	if err != nil {
		return
	}
	w := io.Writer(os.Stdout)
	err = startInstance(w, "confident-sweep-395418", "us-west4-b", "instance-1")
	if err != nil {
		panic(err)
	}
}
