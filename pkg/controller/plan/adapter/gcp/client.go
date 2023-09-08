package gcp

import (
	"errors"
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/ref"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	gcpclient "github.com/konveyor/forklift-controller/pkg/lib/client/gcp"
	liberr "github.com/konveyor/forklift-controller/pkg/lib/error"
)

var ResourceNotFoundError = errors.New("resource not found")
var NameOrIDRequiredError = errors.New("id or name is required")
var UnexpectedVolumeStatusError = errors.New("unexpected volume status")

type Client struct {
	gcpclient.Client
	Context *plancontext.Context
}

func (r *Client) PreTransferActions(vmRef ref.Ref) (ready bool, err error) {
	//_, err = r.getVM(vmRef)
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

// Power on the source VM.
func (r *Client) PowerOn(vmRef ref.Ref) (err error) {
	err = r.VMStart(vmRef.Namespace, vmRef.Type, vmRef.Name)
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

//func (r *Client) getVM(vmRef ref.Ref) (vm *computepb.Instance, err error) {
//	if vmRef.Name == "" && vmRef.ID == "" {
//		err = NameOrIDRequiredError
//		return
//	}
//	// convert string vmRef.ID to uint64
//	vmID, err := strconv.ParseUint(vmRef.ID, 10, 64)
//	if err != nil {
//		err = liberr.Wrap(err)
//	}
//	err = c.connectComputeServiceAPI()
//	req := &computepb.AggregatedListInstancesRequest{
//		Project:    vmRef.Namespace,
//		MaxResults: proto.Uint32(3),
//	}
//	it := r.computeService.AggregatedList(r.ctx, req)
//	for {
//		pair, err := it.Next()
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			return
//		}
//		instances := pair.Value.Instances
//		if len(instances) > 0 {
//			for _, instance := range instances {
//				if instance.GetName() == vmRef.Name || instance.GetId() == vmID {
//					vm = instance
//					return
//				}
//			}
//		}
//	}
//	err = ResourceNotFoundError
//	return
//}
