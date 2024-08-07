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

package mockbigqueryconnection

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/common/projects"
	pb "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/bigquery/connection/v1"
)

type ConnectionV1 struct {
	*MockService
	pb.UnimplementedConnectionServiceServer
}

func (s *ConnectionV1) GetConnection(ctx context.Context, req *pb.GetConnectionRequest) (*pb.Connection, error) {
	name, err := s.parseConnectionName(req.Name)
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	obj := &pb.Connection{}
	if err := s.storage.Get(ctx, fqn, obj); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "Not found: Connection %s", fqn)
		}
		return nil, err
	}

	return obj, nil
}

func (s *ConnectionV1) CreateConnection(ctx context.Context, req *pb.CreateConnectionRequest) (*pb.Connection, error) {
	reqName := req.Parent + "/connections/" + req.ConnectionId
	name, err := s.parseConnectionName(reqName)
	if err != nil {
		return nil, err
	}

	fqn := name.String()
	now := time.Now()

	obj := proto.Clone(req.Connection).(*pb.Connection)

	obj.Name = name.stringInResponse()
	obj.CreationTime = now.Unix()
	obj.LastModifiedTime = now.Unix()
	if obj.GetCloudResource().GetServiceAccountId() == "" {
		obj.GetCloudResource().ServiceAccountId = "bqcx-${projectNumber}-abcd@gcp-sa-bigquery-condel.iam.gserviceaccount.com"
	}

	if err := s.storage.Create(ctx, fqn, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *ConnectionV1) UpdateConnection(ctx context.Context, req *pb.UpdateConnectionRequest) (*pb.Connection, error) {
	name, err := s.parseConnectionName(req.Name)
	if err != nil {
		return nil, err
	}

	fqn := name.String()
	now := time.Now()

	obj := &pb.Connection{}
	if err := s.storage.Get(ctx, fqn, obj); err != nil {
		return nil, err
	}

	paths := req.GetUpdateMask().GetPaths()
	for _, path := range paths {
		switch path {
		case "friendlyName":
			obj.FriendlyName = req.GetConnection().GetFriendlyName()
		case "description":
			obj.Description = req.GetConnection().GetDescription()
		default:
			return nil, status.Errorf(codes.InvalidArgument, "update_mask path %q not valid", path)
		}
	}
	obj.LastModifiedTime = now.Unix()
	if err := s.storage.Update(ctx, fqn, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *ConnectionV1) DeleteConnection(ctx context.Context, req *pb.DeleteConnectionRequest) (*empty.Empty, error) {
	name, err := s.parseConnectionName(req.Name)
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	oldObj := &pb.Connection{}
	if err := s.storage.Delete(ctx, fqn, oldObj); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

type connectionName struct {
	Project    *projects.ProjectData
	Location   string
	ResourceID string
}

func (n *connectionName) String() string {
	return "projects/" + n.Project.ID + "/locations/" + n.Location + "/connections/" + n.ResourceID
}

func (n *connectionName) stringInResponse() string {
	// The returned "name" in the response uses the project number instead of
	// project ID.
	return "projects/${projectNumber}/locations/" + n.Location + "/connections/" + n.ResourceID
}

// parseConnectionName parses a string into a connectionName.
// The expected form is projects/<projectID>/locations/<location>/connections/<connectionID>
func (s *MockService) parseConnectionName(name string) (*connectionName, error) {
	tokens := strings.Split(name, "/")

	if len(tokens) == 6 && tokens[0] == "projects" && tokens[2] == "locations" && tokens[4] == "connections" {
		project, err := s.Projects.GetProjectByID(tokens[1])
		if err != nil {
			return nil, err
		}

		name := &connectionName{
			Project:    project,
			Location:   tokens[3],
			ResourceID: tokens[5],
		}

		return name, nil
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "name %q is not valid", name)
	}
}
