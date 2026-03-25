// Copyright 2026 Google LLC
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

package clouddeploy

import (
	"context"
	"fmt"

	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/clouddeploy/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/common"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/label"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/structuredreporting"

	gcp "cloud.google.com/go/deploy/apiv1"
	clouddeploypb "cloud.google.com/go/deploy/apiv1/deploypb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

func init() {
	registry.RegisterModel(krm.CloudDeployTargetGVK, NewTargetModel)
}

func NewTargetModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &modelTarget{config: *config}, nil
}

var _ directbase.Model = &modelTarget{}

type modelTarget struct {
	config config.ControllerConfig
}

func (m *modelTarget) client(ctx context.Context) (*gcp.CloudDeployClient, error) {
	var opts []option.ClientOption
	opts, err := m.config.RESTClientOptions()
	if err != nil {
		return nil, err
	}
	gcpClient, err := gcp.NewCloudDeployRESTClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("building Target client: %w", err)
	}
	return gcpClient, err
}

func (m *modelTarget) AdapterForObject(ctx context.Context, op *directbase.AdapterForObjectOperation) (directbase.Adapter, error) {
	u := op.GetUnstructured()
	reader := op.Reader
	obj := &krm.CloudDeployTarget{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	id, err := obj.GetIdentity(ctx, reader)
	if err != nil {
		return nil, err
	}

	// Get clouddeploy GCP client
	gcpClient, err := m.client(ctx)
	if err != nil {
		return nil, err
	}

	mapCtx := &direct.MapContext{}
	desiredPb := CloudDeployTargetSpec_ToProto(mapCtx, &obj.Spec)
	if mapCtx.Err() != nil {
		return nil, mapCtx.Err()
	}
	desiredPb.Labels = label.NewGCPLabelsFromK8sLabels(u.GetLabels())

	return &TargetAdapter{
		id:        id.(*krm.CloudDeployTargetIdentity),
		gcpClient: gcpClient,
		desiredPb: desiredPb,
	}, nil
}

func (m *modelTarget) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	id := &krm.CloudDeployTargetIdentity{}
	if err := id.FromExternal(url); err != nil {
		return nil, err
	}

	// Get clouddeploy GCP client
	gcpClient, err := m.client(ctx)
	if err != nil {
		return nil, err
	}
	return &TargetAdapter{
		id:        id,
		gcpClient: gcpClient,
	}, nil
}

type TargetAdapter struct {
	id        *krm.CloudDeployTargetIdentity
	gcpClient *gcp.CloudDeployClient
	desiredPb *clouddeploypb.Target
	actual    *clouddeploypb.Target
}

var _ directbase.Adapter = &TargetAdapter{}

// Find retrieves the GCP resource.
// Return true means the object is found. This triggers Adapter `Update` call.
// Return false means the object is not found. This triggers Adapter `Create` call.
// Return a non-nil error requeues the requests.
func (a *TargetAdapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("getting Target", "name", a.id)

	req := &clouddeploypb.GetTargetRequest{Name: a.id.String()}
	targetpb, err := a.gcpClient.GetTarget(ctx, req)
	if err != nil {
		if direct.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("getting Target %q: %w", a.id, err)
	}

	a.actual = targetpb
	return true, nil
}

// Create creates the resource in GCP based on `spec` and update the Config Connector object `status` based on the GCP response.
func (a *TargetAdapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("creating Target", "name", a.id)
	mapCtx := &direct.MapContext{}

	req := &clouddeploypb.CreateTargetRequest{
		Parent:   a.id.Parent().String(),
		Target:   a.desiredPb,
		TargetId: a.id.ID(),
	}
	op, err := a.gcpClient.CreateTarget(ctx, req)
	if err != nil {
		return fmt.Errorf("creating Target %s: %w", a.id, err)
	}
	created, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Target %s waiting creation: %w", a.id, err)
	}
	log.V(2).Info("successfully created Target", "name", a.id)

	status := &krm.CloudDeployTargetStatus{}
	status.ObservedState = CloudDeployTargetObservedState_FromProto(mapCtx, created)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	status.ExternalRef = direct.LazyPtr(a.id.String())
	return createOp.UpdateStatus(ctx, status, nil)
}

// Update updates the resource in GCP based on `spec` and update the Config Connector object `status` based on the GCP response.
func (a *TargetAdapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("updating Target", "name", a.id)
	mapCtx := &direct.MapContext{}

	a.desiredPb.Name = a.id.String()

	paths := make(sets.Set[string])
	var err error

	// etag and executionConfigs are server-generated and not in KCC spec.
	// We zero them out to avoid unnecessary drift.
	// TODO: We should probably handle executionConfigs better if they are actually in the spec
	a.actual.Etag = ""
	a.actual.ExecutionConfigs = nil

	paths, err = common.CompareProtoMessage(a.desiredPb, a.actual, common.BasicDiff)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		log.V(2).Info("no field needs update", "name", a.id)
		return nil
	}
	report := &structuredreporting.Diff{Object: updateOp.GetUnstructured()}
	for path := range paths {
		report.AddField(path, nil, nil)
	}
	structuredreporting.ReportDiff(ctx, report)

	updateMask := &fieldmaskpb.FieldMask{
		Paths: sets.List(paths),
	}

	req := &clouddeploypb.UpdateTargetRequest{
		UpdateMask: updateMask,
		Target:     a.desiredPb,
	}
	op, err := a.gcpClient.UpdateTarget(ctx, req)
	if err != nil {
		return fmt.Errorf("updating Target %s: %w", a.id, err)
	}
	updated, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Target %s waiting update: %w", a.id, err)
	}
	log.V(2).Info("successfully updated Target", "name", a.id)

	status := &krm.CloudDeployTargetStatus{}
	status.ObservedState = CloudDeployTargetObservedState_FromProto(mapCtx, updated)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	return updateOp.UpdateStatus(ctx, status, nil)
}

// Export maps the GCP object to a Config Connector resource `spec`.
func (a *TargetAdapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	if a.actual == nil {
		return nil, fmt.Errorf("Find() not called")
	}
	u := &unstructured.Unstructured{}

	obj := &krm.CloudDeployTarget{}
	mapCtx := &direct.MapContext{}
	obj.Spec = direct.ValueOf(CloudDeployTargetSpec_FromProto(mapCtx, a.actual))
	if mapCtx.Err() != nil {
		return nil, mapCtx.Err()
	}
	// TODO: Handle Labels in Export if needed (it is usually handled by the generic exporter)

	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	u.SetName(a.id.ID())
	u.SetGroupVersionKind(krm.CloudDeployTargetGVK)

	u.Object = uObj
	return u, nil
}

// Delete the resource from GCP service when the corresponding Config Connector resource is deleted.
func (a *TargetAdapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("deleting Target", "name", a.id)

	req := &clouddeploypb.DeleteTargetRequest{Name: a.id.String()}
	op, err := a.gcpClient.DeleteTarget(ctx, req)
	if err != nil {
		if direct.IsNotFound(err) {
			// Return success if not found (assume it was already deleted).
			log.V(2).Info("skipping delete for non-existent Target, assuming it was already deleted", "name", a.id)
			return true, nil
		}
		return false, fmt.Errorf("deleting Target %s: %w", a.id, err)
	}
	log.V(2).Info("successfully deleted Target", "name", a.id)

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("waiting delete Target %s: %w", a.id, err)
	}
	return true, nil
}
