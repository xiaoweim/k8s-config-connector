// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquerydataset

import (
	"context"
	"fmt"
	"reflect"

	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/bigquery/v1beta1"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/bigquery/v2"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"

	bigquerypb "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/bigquery/v2"
	api "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ctrlName      = "bigquery-controller"
	serviceDomain = "//bigquery.googleapis.com"
)

func init() {
	registry.RegisterModel(krm.BigQueryDatasetGVK, NewModel)
}

func NewModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &model{config: *config}, nil
}

var _ directbase.Model = &model{}

type model struct {
	config config.ControllerConfig
}

func (m *model) service(ctx context.Context) (*api.Service, error) {
	var opts []option.ClientOption
	opts, err := m.config.RESTClientOptions()
	if err != nil {
		return nil, err
	}
	gcpService, err := api.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("building Dataset client: %w", err)
	}
	return gcpService, err
}

func (m *model) AdapterForObject(ctx context.Context, reader client.Reader, u *unstructured.Unstructured) (directbase.Adapter, error) {
	obj := &krm.BigQueryDataset{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	// Get ResourceID
	resourceID := direct.ValueOf(obj.Spec.ResourceID)
	if resourceID == "" {
		resourceID = obj.GetName()
	}
	if resourceID == "" {
		return nil, fmt.Errorf("cannot resolve resource ID")
	}

	projectRef, err := refs.ResolveProject(ctx, reader, obj, obj.Spec.ProjectRef)
	if err != nil {
		return nil, err
	}
	projectID := projectRef.ProjectID
	if projectID == "" {
		return nil, fmt.Errorf("cannot resolve project")
	}

	// Resolve KMS key reference
	if obj.Spec.DefaultEncryptionConfiguration != nil {
		kmsRef, err := refs.ResolveKMSCryptoKeyRef(ctx, reader, obj, obj.Spec.DefaultEncryptionConfiguration.KmsKeyRef)
		if err != nil {
			return nil, err
		}
		obj.Spec.DefaultEncryptionConfiguration.KmsKeyRef = kmsRef
	}

	id, err := krm.NewBigQueryDatasetRef(ctx, reader, obj)
	if err != nil {
		return nil, err
	}

	// Get bigquery GCP client
	gcpService, err := m.service(ctx)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		id:         id,
		gcpService: gcpService,
		desired:    obj,
	}, nil
}

func (m *model) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	// TODO: Support URLs
	return nil, nil
}

type Adapter struct {
	id         *krm.BigQueryDatasetRef
	gcpService *api.Service
	desired    *krm.BigQueryDataset
	actual     *bigquerypb.Dataset
}

var _ directbase.Adapter = &Adapter{}

func (a *Adapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx).WithName(ctrlName)
	log.V(2).Info("getting BigQueryDataset", "name", a.id.External)

	parent, datasetId, err := krm.ParseBigQueryDatasetExternal(a.id.External)
	if err != nil {
		return false, fmt.Errorf("failed to parse bigquery dataset full name, %w", err)
	}
	req := &bigquerypb.GetDatasetRequest{
		DatasetId: &datasetId,
		ProjectId: &parent.ProjectID,
	}
	datasetGetCall := a.gcpService.Datasets.Get(*req.ProjectId, *req.DatasetId)
	datasetpb, err := datasetGetCall.Do()
	if err != nil {
		if direct.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("getting BigQueryDataset %q: %w", a.id.External, err)
	}
	if err := convertAPIToProto(datasetpb, &a.actual); err != nil {
		return false, nil
	}
	return true, nil
}

func (a *Adapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {

	log := klog.FromContext(ctx).WithName(ctrlName)
	log.V(2).Info("creating Dataset", "name", a.id.External)
	mapCtx := &direct.MapContext{}

	desiredDataset := BigQueryDataset_CovertKRMToAPI(mapCtx, a.desired)
	desiredDataset.Labels = make(map[string]string)
	for k, v := range a.desired.GetObjectMeta().GetLabels() {
		desiredDataset.Labels[k] = v
	}
	desiredDataset.Labels["managed-by-cnrm"] = "true"
	parent, _, err := krm.ParseBigQueryDatasetExternal(a.id.External)
	if err != nil {
		return fmt.Errorf("failed to parse bigquery dataset full name, %w", err)
	}
	insertDatasetCall := a.gcpService.Datasets.Insert(parent.ProjectID, desiredDataset)
	inserted, err := insertDatasetCall.Do()
	if err != nil {
		return fmt.Errorf("inserting Dataset %s: %w", a.id.External, err)
	}
	log.V(2).Info("successfully inserted Dataset", "name", a.id.External)

	status := &krm.BigQueryDatasetStatus{}
	status = BigQueryDatasetStatusObservedState_FromProto(mapCtx, inserted)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	return createOp.UpdateStatus(ctx, status, nil)
}

func (a *Adapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	u := updateOp.GetUnstructured()

	log := klog.FromContext(ctx).WithName(ctrlName)
	log.V(2).Info("updating Dataset", "name", a.id.External)
	mapCtx := &direct.MapContext{}

	// Convert KRM object to proto message
	desiredKRM := a.desired.DeepCopy()
	desired := BigQueryDatasetSpec_ToProto(mapCtx, &desiredKRM.Spec, desiredKRM.Name)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	actual := a.actual
	resource := proto.Clone(a.actual).(*bigquerypb.Dataset) // this is the proto resource we are passing to GCP API update call.

	// Check for immutable fields
	if !reflect.DeepEqual(desired.Location, *actual.Location) {
		return fmt.Errorf("BigQueryDataset %s/%s location cannot be changed, actual: %s, desired: %s", u.GetNamespace(), u.GetName(), actual.GetLocation(), desired.Location)
	}

	// Find diff
	updateMask := &fieldmaskpb.FieldMask{}
	if !reflect.DeepEqual(desired.Description, actual.GetDescription()) {
		resource.Description = direct.LazyPtr(desired.Description)
		updateMask.Paths = append(updateMask.Paths, "description")
	}
	if !reflect.DeepEqual(desired.FriendlyName, actual.GetFriendlyName()) {
		resource.FriendlyName = direct.LazyPtr(desired.FriendlyName)
		updateMask.Paths = append(updateMask.Paths, "friendly_name")
	}
	if !reflect.DeepEqual(desired.DefaultPartitionExpirationMs, actual.GetDefaultPartitionExpirationMs()) {
		resource.DefaultPartitionExpirationMs = direct.LazyPtr(desired.DefaultPartitionExpirationMs)
		updateMask.Paths = append(updateMask.Paths, "default_partition_expirationMs")
	}
	if !reflect.DeepEqual(desired.DefaultTableExpirationMs, actual.GetDefaultTableExpirationMs()) {
		resource.DefaultTableExpirationMs = direct.LazyPtr(desired.DefaultTableExpirationMs)
		updateMask.Paths = append(updateMask.Paths, "default_table_expirationMs")
	}
	if !reflect.DeepEqual(desired.DefaultCollation, actual.GetDefaultCollation()) {
		resource.DefaultCollation = direct.LazyPtr(desired.DefaultCollation)
		updateMask.Paths = append(updateMask.Paths, "default_collation")
	}
	if desired.DefaultEncryptionConfiguration != nil && resource.DefaultEncryptionConfiguration != nil && !reflect.DeepEqual(desired.DefaultEncryptionConfiguration, actual.DefaultEncryptionConfiguration) {
		resource.DefaultEncryptionConfiguration.KmsKeyName = direct.LazyPtr(desired.DefaultEncryptionConfiguration.KmsKeyName)
		updateMask.Paths = append(updateMask.Paths, "default_encryption_configuration")
	}
	if !reflect.DeepEqual(desired.IsCaseInsensitive, actual.GetIsCaseInsensitive()) {
		resource.IsCaseInsensitive = direct.LazyPtr(desired.IsCaseInsensitive)
		updateMask.Paths = append(updateMask.Paths, "is_case_sensitive")
	}
	if !reflect.DeepEqual(desired.MaxTimeTravelHours, actual.GetMaxTimeTravelHours()) {
		resource.MaxTimeTravelHours = direct.LazyPtr(desired.MaxTimeTravelHours)
		updateMask.Paths = append(updateMask.Paths, "max_time_interval_hours")
	}
	if desired.Access != nil && actual.Access != nil && len(desired.Access) > 0 && !reflect.DeepEqual(desired.Access, actual.Access) {
		for _, access := range desired.Access {
			resource.Access = append(resource.Access, bigquery.Access_ToProto(mapCtx, Access_convertAPIToProto(mapCtx, access)))
		}
		updateMask.Paths = append(updateMask.Paths, "access")
	}
	if !reflect.DeepEqual(desired.StorageBillingModel, actual.GetStorageBillingModel()) {
		resource.StorageBillingModel = direct.LazyPtr(desired.StorageBillingModel)
		updateMask.Paths = append(updateMask.Paths, "storage_billing_model")
	}

	if len(updateMask.Paths) == 0 {
		return nil
	}
	parent, datasetId, err := krm.ParseBigQueryDatasetExternal(a.id.External)
	if err != nil {
		return fmt.Errorf("failed to parse bigquery dataset full name, %w", err)
	}
	req := &bigquerypb.UpdateDatasetRequest{
		DatasetId: &datasetId,
		ProjectId: &parent.ProjectID,
	}
	datasetToUpdate := &api.Dataset{}
	if err := convertProtoToAPI(resource, datasetToUpdate); err != nil {
		return fmt.Errorf("Failed to convert bigquery proto to golang api proto, err: %w", err)
	}
	// BigQueryDataset GCP service adds default access if the access is unset.
	if desired.Access == nil || len(desired.Access) == 0 {
		datasetGetCall := a.gcpService.Datasets.Get(*req.ProjectId, *req.DatasetId)
		existingDataset, err := datasetGetCall.Do()
		if err != nil {
			if direct.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("getting BigQueryDataset %q: %w", a.id.External, err)
		}
		datasetToUpdate.Access = existingDataset.Access
	}
	updateDatasetCall := a.gcpService.Datasets.Update(*req.ProjectId, *req.DatasetId, datasetToUpdate)
	updated, err := updateDatasetCall.Do()
	if err != nil {
		return fmt.Errorf("updating Dataset %s: %w", a.id.External, err)
	}
	log.V(2).Info("successfully updated Dataset", "name", a.id.External)

	status := &krm.BigQueryDatasetStatus{}
	status = BigQueryDatasetStatusObservedState_FromProto(mapCtx, updated)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	return updateOp.UpdateStatus(ctx, status, nil)
}

func (a *Adapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	if a.actual == nil {
		return nil, fmt.Errorf("Find() not called")
	}
	u := &unstructured.Unstructured{}

	obj := &krm.BigQueryDataset{}
	mapCtx := &direct.MapContext{}
	obj.Spec = direct.ValueOf(ConvertKRMDataset_To_BigQueryDatasetSpec(mapCtx, a.actual))
	if mapCtx.Err() != nil {
		return nil, mapCtx.Err()
	}
	parent, _, err := krm.ParseBigQueryDatasetExternal(a.id.External)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bigquery dataset full name, %w", err)
	}

	obj.Spec.ProjectRef = &refs.ProjectRef{Name: parent.ProjectID}
	obj.Spec.Location = &parent.Location
	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	u.Object = uObj
	return u, nil
}

// Delete implements the Adapter interface.
func (a *Adapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	log := klog.FromContext(ctx).WithName(ctrlName)
	log.V(2).Info("deleting Dataset", "name", a.id.External)

	parent, datasetId, err := krm.ParseBigQueryDatasetExternal(a.id.External)
	if err != nil {
		return false, fmt.Errorf("failed to parse bigquery dataset full name, %w", err)
	}
	req := &bigquerypb.DeleteDatasetRequest{
		DatasetId: &datasetId,
		ProjectId: &parent.ProjectID,
	}
	deleteDatasetCall := a.gcpService.Datasets.Delete(*req.ProjectId, *req.DatasetId)
	err = deleteDatasetCall.Do()
	if err != nil {
		return false, fmt.Errorf("deleting Dataset %s: %w", a.id.External, err)
	}
	log.V(2).Info("successfully deleted Dataset", "name", a.id.External)

	return true, nil
}