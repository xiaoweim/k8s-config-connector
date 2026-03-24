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

package main

import (
	"testing"
)

// baseCRD is a minimal CRD for testing.
const baseCRD = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            description: "The spec"
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
                    description: "Project name"
              region:
                type: string
                description: "Region"
          status:
            type: object
            description: "The status"
            properties:
              observedGeneration:
                type: integer
                format: int64
                description: "Observed generation"
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`

func TestEquivalence(t *testing.T) {
	cases := []struct {
		name          string
		old           string
		new           string
		expectedDiffs []string
		expectedNotes []string
	}{
		{
			name: "DescriptionChange",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            description: "The spec - updated description"
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
                    description: "Project name - updated"
              region:
                type: string
                description: "Region - updated"
          status:
            type: object
            description: "The status - updated"
            properties:
              observedGeneration:
                type: integer
                format: int64
                description: "Observed generation - updated"
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: nil,
			expectedNotes: nil,
		},
		{
			name: "AddListKind",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    listKind: FooList
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: nil,
			expectedNotes: []string{`spec.names.listKind added: "FooList" (allowed)`},
		},
		{
			name: "AllowedExternalRef",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
              externalRef:
                type: string
`,
			expectedDiffs: nil,
			expectedNotes: []string{"[v1alpha1] field added under status: status.externalRef (type: string, allowed)"},
		},
		{
			name: "AddSpecField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
              newSpecField:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: []string{"[v1alpha1] field added: spec.newSpecField (type: string)"},
			expectedNotes: nil,
		},
		{
			name: "RemoveField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: []string{"[v1alpha1] field removed: spec.region (was string)"},
			expectedNotes: nil,
		},
		{
			name: "RemoveStatusField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
`,
			expectedDiffs: []string{"[v1alpha1] field removed: status.conditions (was array)"},
			expectedNotes: nil,
		},
		{
			name: "TypeChange",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: integer
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: []string{"[v1alpha1] field type changed: spec.region (string -> integer)"},
			expectedNotes: nil,
		},
		{
			name: "IntegerTypeChange",
			old: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1beta1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              httpKeepAliveTimeoutSec:
                type: integer
              count:
                type: integer
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
              proxyId:
                type: integer
              nodePort:
                type: integer
`,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1beta1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              httpKeepAliveTimeoutSec:
                type: integer
                format: int32
              count:
                type: integer
                format: int64
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              proxyId:
                type: integer
                format: int64
              nodePort:
                type: integer
                format: int32
`,
			expectedDiffs: []string{"[v1beta1] field type changed: spec.count (integer -> int64)"},
			expectedNotes: []string{
				"[v1beta1] field type changed: spec.httpKeepAliveTimeoutSec (integer -> int32) (allowed)",
				"[v1beta1] field type changed: status.nodePort (integer -> int32) (allowed)",
				"[v1beta1] field type changed: status.observedGeneration (integer -> int64) (allowed)",
				"[v1beta1] field type changed: status.proxyId (integer -> int64) (allowed)",
			},
		},
		{
			name: "NewStatusField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
                    newField:
                      type: string
              externalRef:
                type: string
              externalRefFoo:
                type: string
              externalRefName:
                type: string
              observedState:
                type: object
                properties:
                  bar:
                    type: string
              observedStateFoo:
                type: string
              unallowedObject:
                type: object
                properties:
                  child:
                    type: string
`,
			expectedDiffs: []string{
				"[v1alpha1] field added under status: status.conditions[].newField (type: string)",
				"[v1alpha1] field added under status: status.externalRefFoo (type: string)",
				"[v1alpha1] field added under status: status.externalRefName (type: string)",
				"[v1alpha1] field added under status: status.observedStateFoo (type: string)",
				"[v1alpha1] field added under status: status.unallowedObject (type: object)",
			},
			expectedNotes: []string{
				"[v1alpha1] field added under status: status.externalRef (type: string, allowed)",
				"[v1alpha1] field added under status: status.observedState (type: object, allowed)",
				"[v1alpha1] field added under status: status.observedState.bar (type: string, allowed)",
			},
		},
		{
			name: "EmptyObservedState",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
              observedState:
                type: object
`,
			expectedDiffs: nil,
			expectedNotes: []string{"[v1alpha1] field added under status: status.observedState (type: object, allowed)"},
		},
		{
			name: "DisallowedStatusFieldTypeChange",
			old: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          status:
            type: object
            properties:
              externalRef:
                type: string
`,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          status:
            type: object
            properties:
              externalRef:
                type: integer
`,
			expectedDiffs: []string{"[v1alpha1] field type changed: status.externalRef (string -> integer)"},
			expectedNotes: nil,
		},
		{
			name: "ChangeListKind",
			old: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    listKind: OldList
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
`,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    listKind: NewList
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
`,
			expectedDiffs: []string{`spec.names.listKind changed: "OldList" -> "NewList"`},
			expectedNotes: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oldCRD, err := parseCRD([]byte(tc.old))
			if err != nil {
				t.Fatalf("parseCRD old: %v", err)
			}
			newCRD, err := parseCRD([]byte(tc.new))
			if err != nil {
				t.Fatalf("parseCRD new: %v", err)
			}

			result := compareEquivalence(oldCRD, newCRD)

			if len(result.Diffs) != len(tc.expectedDiffs) {
				t.Fatalf("expected %d diffs, got %d: %v", len(tc.expectedDiffs), len(result.Diffs), result.Diffs)
			}
			for i, d := range result.Diffs {
				if d != tc.expectedDiffs[i] {
					t.Errorf("diff mismatch at %d: expected %q, got %q", i, tc.expectedDiffs[i], d)
				}
			}

			if len(result.Notes) != len(tc.expectedNotes) {
				t.Fatalf("expected %d notes, got %d: %v", len(tc.expectedNotes), len(result.Notes), result.Notes)
			}
			for i, n := range result.Notes {
				if n != tc.expectedNotes[i] {
					t.Errorf("note mismatch at %d: expected %q, got %q", i, tc.expectedNotes[i], n)
				}
			}
		})
	}
}

func TestBackwardCompat(t *testing.T) {
	cases := []struct {
		name          string
		old           string
		new           string
		expectedDiffs []string
		expectedNotes []string
	}{
		{
			name: "AddField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
              region:
                type: string
              newField:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: nil,
			expectedNotes: []string{"[v1alpha1] field added: spec.newField (type: string, allowed)"},
		},
		{
			name: "RemoveField",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: []string{"[v1alpha1] field removed: spec.region (was string)"},
			expectedNotes: nil,
		},
		{
			name: "TypeChange",
			old:  baseCRD,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              projectRef:
                type: object
                properties:
                  name:
                    type: integer
              region:
                type: string
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              conditions:
                type: array
                items:
                  type: object
                  properties:
                    type:
                      type: string
                    status:
                      type: string
`,
			expectedDiffs: []string{"[v1alpha1] field type changed: spec.projectRef.name (string -> integer)"},
			expectedNotes: nil,
		},
		{
			name: "IntegerTypeChange",
			old: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1beta1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              httpKeepAliveTimeoutSec:
                type: integer
              count:
                type: integer
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
              proxyId:
                type: integer
              nodePort:
                type: integer
`,
			new: `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
  versions:
  - name: v1beta1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              httpKeepAliveTimeoutSec:
                type: integer
                format: int32
              count:
                type: integer
                format: int64
          status:
            type: object
            properties:
              observedGeneration:
                type: integer
                format: int64
              proxyId:
                type: integer
                format: int64
              nodePort:
                type: integer
                format: int32
`,
			expectedDiffs: []string{"[v1beta1] field type changed: spec.count (integer -> int64)"},
			expectedNotes: []string{
				"[v1beta1] field type changed: spec.httpKeepAliveTimeoutSec (integer -> int32) (allowed)",
				"[v1beta1] field type changed: status.nodePort (integer -> int32) (allowed)",
				"[v1beta1] field type changed: status.observedGeneration (integer -> int64) (allowed)",
				"[v1beta1] field type changed: status.proxyId (integer -> int64) (allowed)",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oldCRD, err := parseCRD([]byte(tc.old))
			if err != nil {
				t.Fatalf("parseCRD old: %v", err)
			}
			newCRD, err := parseCRD([]byte(tc.new))
			if err != nil {
				t.Fatalf("parseCRD new: %v", err)
			}

			result := compareBackwardCompatibility(oldCRD, newCRD)

			if len(result.Diffs) != len(tc.expectedDiffs) {
				t.Fatalf("expected %d diffs, got %d: %v", len(tc.expectedDiffs), len(result.Diffs), result.Diffs)
			}
			for i, d := range result.Diffs {
				if d != tc.expectedDiffs[i] {
					t.Errorf("diff mismatch at %d: expected %q, got %q", i, tc.expectedDiffs[i], d)
				}
			}

			if len(result.Notes) != len(tc.expectedNotes) {
				t.Fatalf("expected %d notes, got %d: %v", len(tc.expectedNotes), len(result.Notes), result.Notes)
			}
			for i, n := range result.Notes {
				if n != tc.expectedNotes[i] {
					t.Errorf("note mismatch at %d: expected %q, got %q", i, tc.expectedNotes[i], n)
				}
			}
		})
	}
}

func TestGitShow_InvalidRef(t *testing.T) {
	// An invalid ref should return an error, not silently report "file is new".
	_, isNew, err := gitShow("nonexistent-ref-xyz123", "compare_test.go")
	if err == nil {
		t.Fatal("expected an error for an invalid git ref, got nil")
	}
	if isNew {
		t.Fatal("an invalid git ref should return an error, not 'file is new'")
	}
}
