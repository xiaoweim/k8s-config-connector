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

package preview

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/resourceconfig"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGetAlternativeControllerExpectedMap(t *testing.T) {
	testCases := []struct {
		name      string
		configMap resourceconfig.ResourcesControllerMap
		expected  map[schema.GroupKind]k8s.ReconcilerType
	}{
		{
			name: "single config with alternative controller",
			configMap: resourceconfig.ResourcesControllerMap{
				{Group: "foo", Kind: "Bar"}: {
					DefaultController:    k8s.ReconcilerType("tf"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("tf"), k8s.ReconcilerType("direct")},
				},
			},
			expected: map[schema.GroupKind]k8s.ReconcilerType{
				{Group: "foo", Kind: "Bar"}: k8s.ReconcilerType("direct"),
			},
		},
		{
			name: "single config without alternative controller",
			configMap: resourceconfig.ResourcesControllerMap{
				{Group: "foo", Kind: "Bar"}: {
					DefaultController:    k8s.ReconcilerType("direct"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("direct")},
				},
			},
			expected: map[schema.GroupKind]k8s.ReconcilerType{},
		},
		{
			name: "multiple configs mixed",
			configMap: resourceconfig.ResourcesControllerMap{
				{Group: "foo", Kind: "Bar"}: {
					DefaultController:    k8s.ReconcilerType("tf"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("tf"), k8s.ReconcilerType("direct")},
				},
				{Group: "storage", Kind: "StorageBucket"}: {
					DefaultController:    k8s.ReconcilerType("direct"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("direct")},
				},
				{Group: "compute", Kind: "Instance"}: {
					DefaultController:    k8s.ReconcilerType("tf"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("tf"), k8s.ReconcilerType("dcl")},
				},
			},
			expected: map[schema.GroupKind]k8s.ReconcilerType{
				{Group: "foo", Kind: "Bar"}:          k8s.ReconcilerType("direct"),
				{Group: "compute", Kind: "Instance"}: k8s.ReconcilerType("dcl"),
			},
		},
		{
			name: "multiple alternatives picks first non-default",
			configMap: resourceconfig.ResourcesControllerMap{
				{Group: "foo", Kind: "Quad"}: {
					DefaultController:    k8s.ReconcilerType("tf"),
					SupportedControllers: []k8s.ReconcilerType{k8s.ReconcilerType("tf"), k8s.ReconcilerType("direct"), k8s.ReconcilerType("dcl")},
				},
			},
			expected: map[schema.GroupKind]k8s.ReconcilerType{
				{Group: "foo", Kind: "Quad"}: k8s.ReconcilerType("direct"),
			},
		},
		{
			name:      "empty config",
			configMap: resourceconfig.ResourcesControllerMap{},
			expected:  map[schema.GroupKind]k8s.ReconcilerType{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getAlternativeControllerExpectedMap(tc.configMap)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected \n%v\n, got \n%v", tc.expected, actual)
			}
		})
	}
}
