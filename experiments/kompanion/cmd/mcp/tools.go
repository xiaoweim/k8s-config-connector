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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (sc *serverContext) handleGetKCCCRDSchema(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing kind: %v", err)), nil
	}

	crdGVR := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	crds, err := sc.dynamicClient.Resource(crdGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list CRDs: %v", err)), nil
	}

	var targetCRD *unstructured.Unstructured
	for _, crd := range crds.Items {
		crdKind, found, _ := unstructured.NestedString(crd.Object, "spec", "names", "kind")
		if found && crdKind == kind {
			group, _, _ := unstructured.NestedString(crd.Object, "spec", "group")
			if strings.HasSuffix(group, ".cnrm.cloud.google.com") {
				targetCRD = &crd
				break
			}
		}
	}

	if targetCRD == nil {
		return mcp.NewToolResultError(fmt.Sprintf("KCC CRD for kind %s not found", kind)), nil
	}

	versions, found, _ := unstructured.NestedSlice(targetCRD.Object, "spec", "versions")
	if !found || len(versions) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf("no versions found for CRD %s", kind)), nil
	}

	// Find the served version, preferably the latest one
	var schemaObj interface{}
	for _, v := range versions {
		verMap := v.(map[string]interface{})
		served, _ := verMap["served"].(bool)
		if served {
			schemaObj, _, _ = unstructured.NestedFieldNoCopy(verMap, "schema", "openAPIV3Schema")
			// We take the first served version for now
			break
		}
	}

	if schemaObj == nil {
		return mcp.NewToolResultError(fmt.Sprintf("no schema found for CRD %s", kind)), nil
	}

	// Keep it lean: extract only spec and status properties
	properties, found, _ := unstructured.NestedMap(schemaObj.(map[string]interface{}), "properties")
	if found {
		leanSchema := make(map[string]interface{})
		if spec, ok := properties["spec"]; ok {
			leanSchema["spec"] = spec
		}
		if status, ok := properties["status"]; ok {
			leanSchema["status"] = status
		}
		schemaObj = leanSchema
	}

	schemaJSON, err := json.MarshalIndent(schemaObj, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal schema: %v", err)), nil
	}

	return mcp.NewToolResultText(string(schemaJSON)), nil
}

func (sc *serverContext) handleApplyKCCYAML(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	yamlStr, err := request.RequireString("yaml")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing yaml: %v", err)), nil
	}

	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(yamlStr), 4096)
	var results []string

	for {
		var obj unstructured.Unstructured
		if err := decoder.Decode(&obj); err != nil {
			if err == io.EOF {
				break
			}
			return mcp.NewToolResultError(fmt.Sprintf("failed to decode YAML: %v", err)), nil
		}

		if len(obj.Object) == 0 {
			continue
		}

		gvk := obj.GroupVersionKind()
		gvr, err := sc.findGVR(gvk)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to find GVR for %v: %v", gvk, err)), nil
		}

		namespace := obj.GetNamespace()
		name := obj.GetName()

		var res interface{}
		var applyErr error

		// Use Server-Side Apply (Patch with SSA)
		data, err := json.Marshal(obj.Object)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal object: %v", err)), nil
		}

		force := true
		res, applyErr = sc.dynamicClient.Resource(gvr).Namespace(namespace).Patch(ctx, name, types.ApplyPatchType, data, metav1.PatchOptions{
			FieldManager: "kompanion-mcp",
			Force:        &force,
		})

		if applyErr != nil && (apierrors.IsNotFound(applyErr) || strings.Contains(applyErr.Error(), "not found")) {
			res, applyErr = sc.dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, &obj, metav1.CreateOptions{
				FieldManager: "kompanion-mcp",
			})
		}

		if applyErr != nil {
			results = append(results, fmt.Sprintf("Failed to apply %s/%s (%s): %v", namespace, name, gvk.Kind, applyErr))
		} else {
			appliedObj := res.(*unstructured.Unstructured)
			results = append(results, fmt.Sprintf("Successfully applied %s/%s (%s), resourceVersion: %s", namespace, name, gvk.Kind, appliedObj.GetResourceVersion()))
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No resources found in YAML."), nil
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (sc *serverContext) handleDescribeKCCResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing kind: %v", err)), nil
	}
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing namespace: %v", err)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing name: %v", err)), nil
	}

	gvr, err := sc.findGVRByKind(kind)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find GVR for kind %s: %v", kind, err)), nil
	}

	obj, err := sc.dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get resource: %v", err)), nil
	}

	status, found, _ := unstructured.NestedMap(obj.Object, "status")
	if !found {
		return mcp.NewToolResultText(fmt.Sprintf("Resource %s/%s found but has no status field.", namespace, name)), nil
	}

	conditions, _, _ := unstructured.NestedSlice(status, "conditions")
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Resource: %s/%s (%s)\n", namespace, name, kind))
	if len(conditions) > 0 {
		sb.WriteString("Conditions:\n")
		for _, c := range conditions {
			cond := c.(map[string]interface{})
			statusVal, _ := cond["status"].(string)
			reason, _ := cond["reason"].(string)
			message, _ := cond["message"].(string)
			typeVal, _ := cond["type"].(string)
			sb.WriteString(fmt.Sprintf("  - Type: %s, Status: %s, Reason: %s\n", typeVal, statusVal, reason))
			if message != "" {
				sb.WriteString(fmt.Sprintf("    Message: %s\n", message))
			}
		}
	} else {
		sb.WriteString("No conditions found.\n")
	}

	statusJSON, _ := json.MarshalIndent(status, "", "  ")
	sb.WriteString("\nFull Status:\n")
	sb.WriteString(string(statusJSON))

	return mcp.NewToolResultText(sb.String()), nil
}

func (sc *serverContext) handleGetKCCResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing kind: %v", err)), nil
	}
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing namespace: %v", err)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing name: %v", err)), nil
	}

	gvr, err := sc.findGVRByKind(kind)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find GVR for kind %s: %v", kind, err)), nil
	}

	obj, err := sc.dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get resource: %v", err)), nil
	}

	yamlData, err := sigyaml.Marshal(obj.Object)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal resource to YAML: %v", err)), nil
	}

	return mcp.NewToolResultText(string(yamlData)), nil
}

func (sc *serverContext) handleDeleteKCCResource(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing kind: %v", err)), nil
	}
	namespace, err := request.RequireString("namespace")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing namespace: %v", err)), nil
	}
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Missing name: %v", err)), nil
	}

	gvr, err := sc.findGVRByKind(kind)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find GVR for kind %s: %v", kind, err)), nil
	}

	err = sc.dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete resource: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted %s/%s (%s)", namespace, name, kind)), nil
}

func (sc *serverContext) handleListKCCKinds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	crdGVR := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	crds, err := sc.dynamicClient.Resource(crdGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list CRDs: %v", err)), nil
	}

	var kinds []string
	for _, crd := range crds.Items {
		group, found, _ := unstructured.NestedString(crd.Object, "spec", "group")
		if found && strings.HasSuffix(group, ".cnrm.cloud.google.com") {
			kind, _, _ := unstructured.NestedString(crd.Object, "spec", "names", "kind")
			kinds = append(kinds, fmt.Sprintf("- %s (%s)", kind, group))
		}
	}

	if len(kinds) == 0 {
		return mcp.NewToolResultText("No KCC CRDs found."), nil
	}

	return mcp.NewToolResultText("Available KCC Kinds:\n" + strings.Join(kinds, "\n")), nil
}

func (sc *serverContext) handleListKCCResources(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kind := request.GetString("kind", "")
	namespace := request.GetString("namespace", "")
	limit := request.GetInt("limit", 100)
	if limit <= 0 {
		limit = 100
	}

	var gvrs []schema.GroupVersionResource
	if kind != "" {
		gvr, err := sc.findGVRByKind(kind)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to find GVR for kind %s: %v", kind, err)), nil
		}
		gvrs = append(gvrs, gvr)
	} else {
		// List all KCC resources
		apiResourceLists, err := sc.discoveryClient.ServerPreferredResources()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get preferred resources: %v", err)), nil
		}
		for _, apiResourceList := range apiResourceLists {
			if !strings.Contains(apiResourceList.GroupVersion, ".cnrm.cloud.google.com/") {
				continue
			}
			gv, _ := schema.ParseGroupVersion(apiResourceList.GroupVersion)
			for _, apiResource := range apiResourceList.APIResources {
				if !strings.Contains(apiResource.Name, "/") { // skip subresources
					gvrs = append(gvrs, schema.GroupVersionResource{
						Group:    gv.Group,
						Version:  gv.Version,
						Resource: apiResource.Name,
					})
				}
			}
		}
	}

	var results []string
	count := 0
	for _, gvr := range gvrs {
		if count >= int(limit) {
			break
		}
		list, err := sc.dynamicClient.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			// Skip errors for individual resource types
			continue
		}
		for _, item := range list.Items {
			if count >= int(limit) {
				break
			}
			projectID := item.GetAnnotations()["cnrm.cloud.google.com/project-id"]
			if projectID == "" {
				projectID = "n/a"
			}
			results = append(results, fmt.Sprintf("- Kind: %s, Namespace: %s, Name: %s, ProjectID: %s", item.GetKind(), item.GetNamespace(), item.GetName(), projectID))
			count++
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No KCC resources found."), nil
	}

	output := strings.Join(results, "\n")
	if count >= int(limit) {
		output += fmt.Sprintf("\n\n(Note: results truncated to limit of %d)", int(limit))
	}

	return mcp.NewToolResultText(output), nil
}

func (sc *serverContext) findGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	apiResourceLists, err := sc.discoveryClient.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	for _, apiResourceList := range apiResourceLists {
		if apiResourceList.GroupVersion != gvk.GroupVersion().String() {
			continue
		}
		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == gvk.Kind && !strings.Contains(apiResource.Name, "/") {
				return schema.GroupVersionResource{
					Group:    gvk.Group,
					Version:  gvk.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}
	return schema.GroupVersionResource{}, fmt.Errorf("GVR not found for %v", gvk)
}

func (sc *serverContext) findGVRByKind(kind string) (schema.GroupVersionResource, error) {
	apiResourceLists, err := sc.discoveryClient.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	for _, apiResourceList := range apiResourceLists {
		if !strings.Contains(apiResourceList.GroupVersion, ".cnrm.cloud.google.com/") {
			continue
		}
		gv, _ := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				return schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}
	return schema.GroupVersionResource{}, fmt.Errorf("GVR not found for kind %s", kind)
}
