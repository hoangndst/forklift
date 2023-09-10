package gcp

import "cloud.google.com/go/compute/apiv1/computepb"

// VM Status
const (
	VMStatusRunning = computepb.Instance_RUNNING
	VMStatusStopped = computepb.Instance_STOPPED
)

// Image Status
const (
	ImageStatusReady    = computepb.Image_READY
	ImageStatusFailed   = computepb.Image_FAILED
	ImageStatusPending  = computepb.Image_PENDING
	ImageStatusDeleting = computepb.Image_DELETING
)

// Image Type
const (
	QCOW2 = "qcow2"
	VMDK  = "vmdk"
	VHDX  = "vhdx"
	VPC   = "vpc"
)

type GetOpts struct {
	ID string
}

type VMListOpts struct {
}

type ImageListOpts struct {
}

type ImageCreateOpts struct {
}

type DeleteOpts struct {
}
