package gcp

import (
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/plan"
	"github.com/konveyor/forklift-controller/pkg/apis/forklift/v1beta1/ref"
	plancontext "github.com/konveyor/forklift-controller/pkg/controller/plan/context"
	core "k8s.io/api/core/v1"
	cnv "kubevirt.io/api/core/v1"
	cdi "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

type Builder struct {
	*plancontext.Context
	// MAC addresses already in use on the destination cluster. k=mac, v=vmName
	macConflictsMap map[string]string
}

func (b Builder) Secret(vmRef ref.Ref, in, object *core.Secret) error {
	//TODO implement me
	panic("implement me")
}

func (b Builder) ConfigMap(vmRef ref.Ref, secret *core.Secret, object *core.ConfigMap) error {
	//TODO implement me
	panic("implement me")
}

// Create the destination Kubevirt VM.
func (b Builder) VirtualMachine(vmRef ref.Ref, object *cnv.VirtualMachineSpec, persistentVolumeClaims []core.PersistentVolumeClaim) error {
	vm := &model.Workload{}
}

func (b Builder) DataVolumes(vmRef ref.Ref, secret *core.Secret, configMap *core.ConfigMap, dvTemplate *cdi.DataVolume) (dvs []cdi.DataVolume, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) Tasks(vmRef ref.Ref) ([]*plan.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) TemplateLabels(vmRef ref.Ref) (labels map[string]string, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) ResolveDataVolumeIdentifier(dv *cdi.DataVolume) string {
	//TODO implement me
	panic("implement me")
}

func (b Builder) ResolvePersistentVolumeClaimIdentifier(pvc *core.PersistentVolumeClaim) string {
	//TODO implement me
	panic("implement me")
}

func (b Builder) PodEnvironment(vmRef ref.Ref, sourceSecret *core.Secret) (env []core.EnvVar, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) LunPersistentVolumes(vmRef ref.Ref) (pvs []core.PersistentVolume, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) LunPersistentVolumeClaims(vmRef ref.Ref) (pvcs []core.PersistentVolumeClaim, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) SupportsVolumePopulators() bool {
	//TODO implement me
	panic("implement me")
}

func (b Builder) PopulatorVolumes(vmRef ref.Ref, annotations map[string]string, secretName string) (pvcNames []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) PopulatorTransferredBytes(persistentVolumeClaim *core.PersistentVolumeClaim) (transferredBytes int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) SetPopulatorDataSourceLabels(vmRef ref.Ref, pvcs []core.PersistentVolumeClaim) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b Builder) GetPopulatorTaskName(pvc *core.PersistentVolumeClaim) (taskName string, err error) {
	//TODO implement me
	panic("implement me")
}
