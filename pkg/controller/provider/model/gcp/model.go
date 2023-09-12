package gcp

import (
	"github.com/konveyor/forklift-controller/pkg/controller/provider/model/base"
	libmodel "github.com/konveyor/forklift-controller/pkg/lib/inventory/model"
)

// Errors
var NotFound = libmodel.NotFound

type InvalidRefError = base.InvalidRefError

const (
	MaxDetail = base.MaxDetail
)

// Types
type Model = base.Model
type ListOptions = base.ListOptions
type Concern = base.Concern
type Ref = base.Ref

// Base OpenStack model.
type Base struct {
	// Managed object ID.
	ID string `sql:"pk"`
	// Name
	Name string `sql:"d0,index(name)"`
	// Revision
	Revision int64 `sql:"incremented,d0,index(revision)"`
}

// Get the PK.
func (m *Base) Pk() string {
	return m.ID
}

// String representation.
func (m *Base) String() string {
	return m.ID
}

// GCP Image model.
type Image struct {
	Base
	RevisionValidated int64             `sql:"d0,index(revision_validated)"`
	Architecture      *string           `sql:""`
	ArchiveSizeBytes  *int64            `sql:""`
	CreationTimestamp *string           `sql:""`
	Description       *string           `sql:""`
	DiskSizeGb        *int64            `sql:""`
	Family            *string           `sql:""`
	Id                *uint64           `sql:""`
	Kind              *string           `sql:""`
	LabelFingerprint  *string           `sql:""`
	Labels            map[string]string `sql:""`
	LicenseCodes      []int64           `sql:""`
	Licenses          []string          `sql:""`
	Name              *string           `sql:""`
	SatisfiesPzs      *bool             `sql:""`
	SelfLink          *string           `sql:""`
	SourceDisk        *string           `sql:""`
	SourceDiskId      *string           `sql:""`
	SourceImage       *string           `sql:""`
	SourceImageId     *string           `sql:""`
	SourceSnapshot    *string           `sql:""`
	SourceSnapshotId  *string           `sql:""`
	SourceType        *string           `sql:""`
	Status            *string           `sql:""`
	StorageLocations  []string          `sql:""`
}

type AttachedDisk struct {
	ID string `sql:"d0,fk(volume)"`
}

type VM struct {
	Base
	RevisionValidated int64               `sql:"d0,index(revisionValidated)"`
	PolicyVersion     int                 `sql:"d0,index(policyVersion)" eq:"-"`
	Concerns          []Concern           `sql:"" eq:"-"`
	CpuPlatform       *string             `sql:""`
	CreationTimestamp *string             `sql:""`
	Description       *string             `sql:""`
	Disks             []*AttachedDisk     `sql:""`
	Fingerprint       *string             `sql:""`
	Hostname          *string             `sql:""`
	NetworkInterfaces []*NetworkInterface `sql:""`
	SelfLink          *string             `sql:""`
	Status            *string             `sql:""`
	Zone              *string             `sql:""`
}

type NetworkInterface struct {
	Ipv6AccessType *string `sql:""`
	Ipv6Address    *string `sql:""`
	Name           *string `sql:""`
	Network        *string `sql:""`
	NetworkIP      *string `sql:""`
}

type Network struct {
	Base
	IPv4Range                             *string  `sql:""`
	CreationTimestamp                     *string  `sql:""`
	Description                           *string  `sql:""`
	EnableUlaInternalIpv6                 *bool    `sql:""`
	FirewallPolicy                        *string  `sql:""`
	GatewayIPv4                           *string  `sql:""`
	Id                                    *uint64  `sql:""`
	InternalIpv6Range                     *string  `sql:""`
	NetworkFirewallPolicyEnforcementOrder *string  `sql:""`
	SelfLink                              *string  `sql:""`
	SelfLinkWithId                        *string  `sql:""`
	Subnetworks                           []string `sql:""`
}

// Determine if current revision has been validated.
func (m *VM) Validated() bool {
	return m.RevisionValidated == m.Revision
}
