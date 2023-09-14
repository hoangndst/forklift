package gcp

import (
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/plan"
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/ref"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	model "github.com/konveyor/forklift-controller/pkg/controller/provider/web/gcp"
	liberr "github.com/konveyor/forklift-controller/pkg/lib/error"
	core "k8s.io/api/core/v1"
	cnv "kubevirt.io/api/core/v1"
	cdi "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type Builder struct {
	*plancontext.Context
	// MAC addresses already in use on the destination cluster. k=mac, v=vmName
	macConflictsMap map[string]string
}

// Template labels
const (
	TemplateOSLabel                 = "os.template.kubevirt.io/%s"
	TemplateWorkloadLabel           = "workload.template.kubevirt.io/%s"
	TemplateWorkloadServer          = "server"
	TemplateWorkloadDesktop         = "desktop"
	TemplateWorkloadHighPerformance = "highperformance"
	TemplateFlavorLabel             = "flavor.template.kubevirt.io/%s"
	TemplateFlavorTiny              = "tiny"
	TemplateFlavorSmall             = "small"
	TemplateFlavorMedium            = "medium"
	TemplateFlavorLarge             = "large"
)

// Annotations
const (
	AnnImportDiskId = "cdi.kubevirt.io/storage.import.volumeId"
)

// OS types
const (
	Linux = "linux"
)

// OS Distros
const (
	ArchLinux   = "arch"
	CentOS      = "centos"
	Debian      = "debian"
	Fedora      = "fedora"
	FreeBSD     = "freebsd"
	Gentoo      = "gentoo"
	Mandrake    = "mandrake"
	Mandriva    = "mandriva"
	MES         = "mes"
	MSDOS       = "msdos"
	NetBSD      = "netbsd"
	Netware     = "netware"
	OpenBSD     = "openbsd"
	OpenSolaris = "opensolaris"
	OpenSUSE    = "opensuse"
	RHEL        = "rhel"
	SLED        = "sled"
	Ubuntu      = "ubuntu"
	Windows     = "windows"
)

// Default Operating Systems
const (
	DefaultWindows = "win10"
	DefaultLinux   = "rhel8.1"
	UnknownOS      = "unknown"
)

func (r Builder) Secret(vmRef ref.Ref, in, object *core.Secret) error {
	//TODO implement me
	panic("implement me")
}

func (r Builder) ConfigMap(vmRef ref.Ref, secret *core.Secret, object *core.ConfigMap) error {
	//TODO implement me
	panic("implement me")
}

// Create the destination Kubevirt VM.
func (r Builder) VirtualMachine(vmRef ref.Ref, object *cnv.VirtualMachineSpec, persistentVolumeClaims []core.PersistentVolumeClaim) (err error) {
	vm := &model.Workload{}
	err = b.Source.Inventory.Find(vm, vmRef)
	if err != nil {
		err = liberr.Wrap(
			err,
			"VM lookup failed.",
			"vm",
			vmRef.String())
		return
	}
	r.mapFirmware(vm, vmSpec)
	r.mapResources(vm, vmSpec)
	r.mapHardwareRng(vm, vmSpec)
	r.mapInput(vm, vmSpec)
	r.mapVideo(vm, vmSpec)
	r.mapDisks(vm, persistentVolumeClaims, vmSpec)
	err = r.mapNetworks(vm, vmSpec)
	if err != nil {
		err = liberr.Wrap(
			err,
			"network mapping failed",
			"vm",
			vmRef.String())
		return
	}
	return
}

func (r *Builder) mapNetworks(vm *model.Workload, object *cnv.VirtualMachineSpec) (err error) {
	var kNetworks []cnv.Network
	var kInterfaces []cnv.Interface

}

func (r Builder) DataVolumes(vmRef ref.Ref, secret *core.Secret, configMap *core.ConfigMap, dvTemplate *cdi.DataVolume) (dvs []cdi.DataVolume, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) Tasks(vmRef ref.Ref) ([]*plan.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) TemplateLabels(vmRef ref.Ref) (labels map[string]string, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) ResolveDataVolumeIdentifier(dv *cdi.DataVolume) string {
	//TODO implement me
	panic("implement me")
}

func (r Builder) ResolvePersistentVolumeClaimIdentifier(pvc *core.PersistentVolumeClaim) string {
	//TODO implement me
	panic("implement me")
}

func (r Builder) PodEnvironment(vmRef ref.Ref, sourceSecret *core.Secret) (env []core.EnvVar, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) LunPersistentVolumes(vmRef ref.Ref) (pvs []core.PersistentVolume, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) LunPersistentVolumeClaims(vmRef ref.Ref) (pvcs []core.PersistentVolumeClaim, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) SupportsVolumePopulators() bool {
	//TODO implement me
	panic("implement me")
}

func (r Builder) PopulatorVolumes(vmRef ref.Ref, annotations map[string]string, secretName string) (pvcNames []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) PopulatorTransferredBytes(persistentVolumeClaim *core.PersistentVolumeClaim) (transferredBytes int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) SetPopulatorDataSourceLabels(vmRef ref.Ref, pvcs []core.PersistentVolumeClaim) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r Builder) GetPopulatorTaskName(pvc *core.PersistentVolumeClaim) (taskName string, err error) {
	//TODO implement me
	panic("implement me")
}
