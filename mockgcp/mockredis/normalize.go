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

package mockredis

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/mockgcpregistry"
)

var _ mockgcpregistry.SupportsNormalization = &MockService{}

func (s *MockService) ConfigureVisitor(url string, replacements mockgcpregistry.NormalizingVisitor) {
	replacements.ReplacePath(".status.observedState.uid", "0123456789abcdef")
	replacements.ReplacePath(".status.observedState.pscConnections[].pscConnectionID", "${pscConnectionID}")
	replacements.ReplacePath(".status.observedState.pscConnections[].address", "10.11.12.13")
	replacements.ReplacePath(".status.observedState.discoveryEndpoints[].address", "10.11.12.13")
	replacements.ReplacePath(".pscConnections[].address", "10.11.12.13")
	replacements.ReplacePath(".response.pscConnections[].address", "10.11.12.13")
	replacements.ReplacePath(".discoveryEndpoints[].address", "10.11.12.13")
	replacements.ReplacePath(".response.discoveryEndpoints[].address", "10.11.12.13")
}

func (s *MockService) Previsit(event mockgcpregistry.Event, replacements mockgcpregistry.NormalizingVisitor) {
	// No-op for now
}
