package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	backupdriverv1 "github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/apis/backupdriver/v1"
	"github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/backupdriver"
	backupdriverTypedV1 "github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/generated/clientset/versioned/typed/backupdriver/v1"
	"github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/install"
	pluginUtil "github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/plugin/util"
	"github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/snapshotUtils"
	"github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/utils"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// PVCBackupItemAction is a backup item action plugin for Velero.
type NewPVCBackupItemAction struct {
	Log logrus.FieldLogger
}

// AppliesTo returns information indicating that the PVCBackupItemAction should be invoked to backup PVCs.
func (p *NewPVCBackupItemAction) AppliesTo() (velero.ResourceSelector, error) {
	p.Log.Info("VSphere PVCBackupItemAction AppliesTo")

	return velero.ResourceSelector{
		IncludedResources: []string{"persistentvolumeclaims"},
	}, nil
}

// Execute recognizes PVCs backed by volumes provisioned by vSphere CNS block volumes
func (p *NewPVCBackupItemAction) Execute(item runtime.Unstructured, backup *velerov1.Backup) (runtime.Unstructured, []velero.ResourceIdentifier, error) {
	// Do nothing if volume snapshots have not been requested in this backup
	//if utils.IsSetToFalse(backup.Spec.SnapshotVolumes) {
	ctx := context.Background()
	if backup.Spec.SnapshotVolumes != nil && *backup.Spec.SnapshotVolumes == false {
		p.Log.Infof("Volume snapshotting not requested for backup %s/%s", backup.Namespace, backup.Name)
		return item, nil, nil
	}

	var pvc corev1.PersistentVolumeClaim
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.UnstructuredContent(), &pvc); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	p.Log.Infof("VSphere PVCBackupItemAction for PVC %s/%s started", pvc.Namespace, pvc.Name)
	var err error
	defer func() {
		p.Log.Infof("VSphere PVCBackupItemAction for PVC %s/%s completed with err: %v", pvc.Namespace, pvc.Name, err)
	}()

	// get the velero namespace and the rest config in k8s cluster
	veleroNs, exist := os.LookupEnv("VELERO_NAMESPACE")
	if !exist {
		errMsg := "Failed to lookup the ENV variable for velero namespace"
		p.Log.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		p.Log.Error("Failed to get the rest config in k8s cluster: %v", err)
		return nil, nil, errors.WithStack(err)
	}

	backupdriverClient, err := backupdriverTypedV1.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// Do not claim a backup repository in local mode
	var backupRepositoryName string
	isLocalMode := utils.GetBool(install.DefaultBackupDriverImageLocalMode, false)
	if !isLocalMode {
		p.Log.Info("Claiming backup repository for snapshot")
		bslName := backup.Spec.StorageLocation
		repositoryParameters := make(map[string]string)

		err = utils.RetrieveParamsFromBSL(repositoryParameters, bslName, restConfig, p.Log)
		if err != nil {
			p.Log.Errorf("Failed to translate BSL to repository parameters: %v", err)
			return nil, nil, errors.WithStack(err)
		}

		backupRepositoryName, err = backupdriver.ClaimBackupRepository(ctx, utils.S3RepositoryDriver, repositoryParameters,
			[]string{pvc.Namespace}, veleroNs, backupdriverClient, p.Log)
		if err != nil {
			p.Log.Errorf("Failed to claim backup repository: %v", err)
			return nil, nil, errors.WithStack(err)
		}
	}
	backupRepository := snapshotUtils.NewBackupRepository(backupRepositoryName)

	objectToSnapshot := corev1.TypedLocalObjectReference{
		APIGroup: &corev1.SchemeGroupVersion.Group,
		Kind:     pvc.Kind,
		Name:     pvc.Name,
	}

	p.Log.Info("Creating a Snapshot CR")
	updatedSnapshot, err := snapshotUtils.SnapshotRef(ctx, backupdriverClient, objectToSnapshot, pvc.Namespace, *backupRepository,
		[]backupdriverv1.SnapshotPhase{backupdriverv1.SnapshotPhaseSnapshotted, backupdriverv1.SnapshotPhaseSnapshotFailed}, p.Log)
	if err != nil {
		p.Log.Errorf("Failed to create a Snapshot CR: %v", err)
		return nil, nil, errors.WithStack(err)
	}
	if updatedSnapshot.Status.Phase == backupdriverv1.SnapshotPhaseSnapshotFailed {
		errMsg := fmt.Sprintf("Failed to create a Snapshot CR: Phase=SnapshotFailed, err=%v", updatedSnapshot.Status.Message)
		p.Log.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	// Persist the snapshot blob as an annotation of PVC
	snapshotAnnotation, err := pluginUtil.GetAnnotationFromSnapshot(updatedSnapshot)
	if err != nil {
		p.Log.Errorf("Failed to marshal Snapshot object: %v", err)
		return nil, nil, errors.WithStack(err)
	}
	vals := map[string]string{
		utils.ItemSnapshotLabel: snapshotAnnotation,
	}
	pluginUtil.AddAnnotations(&pvc.ObjectMeta, vals)

	p.Log.Info("Snapshot completed in plugin")

	var additionalItems []velero.ResourceIdentifier

	pvcMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pvc)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &unstructured.Unstructured{Object: pvcMap}, additionalItems, nil
}
