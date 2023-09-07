package gcp

import (
	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/ref"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	liberr "github.com/konveyor/forklift-controller/pkg/lib/error"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
	"os"
	"strconv"
)

var ResourceNotFoundError = errors.New("resource not found")
var NameOrIDRequiredError = errors.New("id or name is required")
var UnexpectedVolumeStatusError = errors.New("unexpected volume status")

type Client struct {
	googleAuth      string
	Context         *plancontext.Context
	ctx             context.Context
	computeMetadata *metadata.Client
	computeService  *compute.InstancesClient
	imageService    *compute.ImagesClient
	networkService  *compute.NetworksClient
	storageService  *storage.Client
}

func (r *Client) Authenticate() (err error) {
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", r.googleAuth)
	if err != nil {
		return
	}
	return
}

func (r *Client) connectComputeMetadataAPI() (err error) {
	if r.computeMetadata == nil {

		if err != nil {
			return
		}
	}
	return
}

func (r *Client) connectComputeServiceAPI() (err error) {
	err = r.Authenticate()
	if err != nil {
		return
	}
	if r.computeService == nil {
		instancesClient, err := compute.NewInstancesRESTClient(r.ctx)
		if err != nil {
			r.computeService = instancesClient
		}
	}
	return
}

func (r *Client) PowerOff(vmRef ref.Ref) error {
	err := r.connectComputeServiceAPI()
	if err != nil {
		return err
	}
	req := &computepb.StopInstanceRequest{
		Project:  vmRef.Namespace,
		Instance: vmRef.Name,
		Zone:     vmRef.Type,
	}
	_, err = r.computeService.Stop(r.ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (r *Client) PowerOn(vmRef ref.Ref) (err error) {
	err = r.connectComputeServiceAPI()
	if err != nil {
		return err
	}
	req := &computepb.StartInstanceRequest{
		Project:  vmRef.Namespace,
		Instance: vmRef.Name,
		Zone:     vmRef.Type,
	}
	_, err = r.computeService.Start(r.ctx, req)
	if err != nil {
		return err
	}
	return
}

func (r *Client) PreTransferActions(vmRef ref.Ref) (ready bool, err error) {
	_, err = r.getVM(vmRef)
	if err != nil {
		err = liberr.Wrap(
			err,
			"VM lookup failed.",
			"vm",
			vmRef.String())
		return
	}
	return
}

func (r *Client) getVM(vmRef ref.Ref) (vm *computepb.Instance, err error) {
	if vmRef.Name == "" && vmRef.ID == "" {
		err = NameOrIDRequiredError
		return
	}
	// convert string vmRef.ID to uint64
	vmID, err := strconv.ParseUint(vmRef.ID, 10, 64)
	if err != nil {
		err = liberr.Wrap(err)
	}
	err = r.connectComputeServiceAPI()
	req := &computepb.AggregatedListInstancesRequest{
		Project:    vmRef.Namespace,
		MaxResults: proto.Uint32(3),
	}
	it := r.computeService.AggregatedList(r.ctx, req)
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
				if instance.GetName() == vmRef.Name || instance.GetId() == vmID {
					vm = instance
					return
				}
			}
		}
	}
	err = ResourceNotFoundError
	return
}
