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

package bigquery

import (
	"testing"

	bigquery "google.golang.org/api/bigquery/v2"
)

func TestPolicyTagsEqual(t *testing.T) {
	tests := []struct {
		name string
		a    *bigquery.TableFieldSchemaPolicyTags
		b    *bigquery.TableFieldSchemaPolicyTags
		want bool
	}{
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "one nil, one empty",
			a:    nil,
			b:    &bigquery.TableFieldSchemaPolicyTags{},
			want: true, // This is what we want to fix
		},
		{
			name: "one empty, one nil",
			a:    &bigquery.TableFieldSchemaPolicyTags{},
			b:    nil,
			want: true, // This is what we want to fix
		},
		{
			name: "both empty",
			a:    &bigquery.TableFieldSchemaPolicyTags{},
			b:    &bigquery.TableFieldSchemaPolicyTags{},
			want: true,
		},
		{
			name: "one with nil names, one with empty names",
			a:    &bigquery.TableFieldSchemaPolicyTags{Names: nil},
			b:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{}},
			want: true,
		},
		{
			name: "one with names, one nil",
			a:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{"tag1"}},
			b:    nil,
			want: false,
		},
		{
			name: "both with same names",
			a:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{"tag1", "tag2"}},
			b:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{"tag2", "tag1"}},
			want: true,
		},
		{
			name: "different names",
			a:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{"tag1"}},
			b:    &bigquery.TableFieldSchemaPolicyTags{Names: []string{"tag2"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := policyTagsEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("policyTagsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoriesEqual(t *testing.T) {
	tests := []struct {
		name string
		a    *bigquery.TableFieldSchemaCategories
		b    *bigquery.TableFieldSchemaCategories
		want bool
	}{
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "one nil, one empty",
			a:    nil,
			b:    &bigquery.TableFieldSchemaCategories{},
			want: true,
		},
		{
			name: "one empty, one nil",
			a:    &bigquery.TableFieldSchemaCategories{},
			b:    nil,
			want: true,
		},
		{
			name: "both empty",
			a:    &bigquery.TableFieldSchemaCategories{},
			b:    &bigquery.TableFieldSchemaCategories{},
			want: true,
		},
		{
			name: "one with nil names, one with empty names",
			a:    &bigquery.TableFieldSchemaCategories{Names: nil},
			b:    &bigquery.TableFieldSchemaCategories{Names: []string{}},
			want: true,
		},
		{
			name: "one with names, one nil",
			a:    &bigquery.TableFieldSchemaCategories{Names: []string{"cat1"}},
			b:    nil,
			want: false,
		},
		{
			name: "both with same names",
			a:    &bigquery.TableFieldSchemaCategories{Names: []string{"cat1", "cat2"}},
			b:    &bigquery.TableFieldSchemaCategories{Names: []string{"cat2", "cat1"}},
			want: true,
		},
		{
			name: "different names",
			a:    &bigquery.TableFieldSchemaCategories{Names: []string{"cat1"}},
			b:    &bigquery.TableFieldSchemaCategories{Names: []string{"cat2"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := categoriesEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("categoriesEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
