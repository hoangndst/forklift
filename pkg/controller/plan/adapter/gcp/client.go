package gcp

import (
	"cloud.google.com/go/compute/apiv1/computepb"
	"errors"
	planapi "github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/plan"
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/ref"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	gcpclient "github.com/konveyor/forklift-controller/pkg/lib/client/gcp"
	liberr "github.com/konveyor/forklift-controller/pkg/lib/error"
	"github.com/konveyor/forklift-controller/pkg/lib/logging"
	cdi "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const (
	ImageStatusReady = "READY"
)

var ResourceNotFoundError = errors.New("resource not found")
var NameOrIDRequiredError = errors.New("id or name is required")
var UnexpectedVolumeStatusError = errors.New("unexpected volume status")

type Client struct {
	gcpclient.Client
	Log     logging.LevelLogger
	Context *plancontext.Context
}

// Return the source VM's power state.
func (r *Client) PowerState(vmRef ref.Ref) (state string, err error) {
	state, err = r.VMStatus(vmRef.Name)
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

// Power on the source VM.
func (r *Client) PowerOn(vmRef ref.Ref) (err error) {
	//vm, err := r.getVM(vmRef)
	err = r.VMStart(vmRef.Name)
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

// Power off the source VM.
func (r *Client) PowerOff(vmRef ref.Ref) (err error) {
	//vm, err := r.getVM(vmRef)
	err = r.VMStop(vmRef.Name)
	if err != nil {
		err = liberr.Wrap(err)
	}
	return
}

// Return whether the source VM is powered off.
func (r *Client) PoweredOff(vmRef ref.Ref) (off bool, err error) {
	state, err := r.PowerState(vmRef)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	off = state == gcpclient.VMStatusStopped.String()
	return
}

// Create a snapshot of the source VM.
func (r *Client) CreateSnapshot(vmRef ref.Ref) (imageID string, err error) {
	return
}

// Remove all warm migration snapshots.
func (r *Client) RemoveSnapshots(vmRef ref.Ref, precopies []planapi.Precopy) (err error) {
	return
}

// Check if a snapshot is ready to transfer.
func (r *Client) CheckSnapshotReady(vmRef ref.Ref, imageID string) (ready bool, err error) {
	return
}

// Set DataVolume checkpoints.
func (r *Client) SetCheckpoints(vmRef ref.Ref, precopies []planapi.Precopy, datavolumes []cdi.DataVolume, final bool) error {
	return nil
}

// Close connections to the provider API.
func (r *Client) Close() {
}

func (r *Client) Finalize(vms []*planapi.VMStatus, planName string) {
	for _, vmStatus := range vms {
		vmRef := ref.Ref{ID: vmStatus.Ref.ID, Name: vmStatus.Ref.Name}
		vm, err := r.getVM(vmRef)
		if err != nil {
			r.Log.Error(err, "failed to find vm", "vm", vm.Name)
			return
		}
		// delete image
	}
}

func (r *Client) DetachDisks(vmRef ref.Ref) (err error) {
	// no-op
	return
}

func (r *Client) getVM(vmRef ref.Ref) (vm *computepb.Instance, err error) {
	if vmRef.Name == "" && vmRef.ID == "" {
		err = NameOrIDRequiredError
		return
	}
	err = r.Client.Get(vm, vmRef.Name)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	return
}

func (r *Client) createImage(vmRef ref.Ref) (imageName string, err error) {
	if vmRef.Name == "" && vmRef.ID == "" {
		err = NameOrIDRequiredError
		return
	}
	// create image from boot disk when instance is stopped
	imageName, err = r.Client.VMCreateImageFromDisk(vmRef.Name, false)
	if err != nil {
		err = liberr.Wrap(err)
		return
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

	// Turn off the VM. (make sure VM is off)
	poweredOff, err := r.PoweredOff(vmRef)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	if !poweredOff {
		err = r.PowerOff(vmRef)
		if err != nil {
			err = liberr.Wrap(err)
			return
		}
	}

	// Create Image from VM's boot disk. (make sure image is ready)
	imageName, err := r.createImage(vmRef)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	image := &computepb.Image{}
	err = r.Client.Get(image, imageName)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	if *image.Status != ImageStatusReady {
		err = liberr.Wrap(ResourceNotFoundError)
		return
	}
	// Export Image to bucket (make sure image in bucket is ready)

	return true, nil
}
