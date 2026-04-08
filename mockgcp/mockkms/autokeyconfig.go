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

// +tool:mockgcp-support
// proto.service: google.cloud.kms.v1.AutokeyAdmin
// proto.message: google.cloud.kms.v1.AutokeyConfig

package mockkms

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/kms/v1"
)

type autokeyAdminServer struct {
	*MockService
	pb.UnimplementedFoldersServerServer
	pb.UnimplementedOrganizationsServerServer
	pb.UnimplementedProjectsServerServer
}

func (r *autokeyAdminServer) GetAutokeyConfigFolder(ctx context.Context, req *pb.GetAutokeyConfigFolderRequest) (*pb.AutokeyConfig, error) {
	name, err := r.parseAutokeyConfigName(req.GetName())
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	obj := &pb.AutokeyConfig{}
	if err := r.storage.Get(ctx, fqn, obj); err != nil {
		if status.Code(err) == codes.NotFound {
			obj.State = pb.AutokeyConfig_UNINITIALIZED.Enum()
			r.storage.Create(ctx, fqn, obj)
			return obj, nil
		}
		return nil, err
	}

	return obj, nil
}

func (r *autokeyAdminServer) UpdateAutokeyConfigFolder(ctx context.Context, req *pb.UpdateAutokeyConfigFolderRequest) (*pb.AutokeyConfig, error) {
	reqName := req.GetFolder().GetName()
	name, err := r.parseAutokeyConfigName(reqName)
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	obj := proto.Clone(req.GetFolder()).(*pb.AutokeyConfig)
	obj.Name = &fqn
	if len(req.GetFolder().GetKeyProject()) > 0 {
		obj.State = pb.AutokeyConfig_ACTIVE.Enum()
	} else {
		obj.State = pb.AutokeyConfig_UNINITIALIZED.Enum()
	}
	if err := r.storage.Update(ctx, fqn, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *autokeyAdminServer) ShowEffectiveAutokeyConfigProject(ctx context.Context, req *pb.ShowEffectiveAutokeyConfigProjectRequest) (*pb.ShowEffectiveAutokeyConfigResponse, error) {
	project := req.GetParent()
	obj := &pb.ShowEffectiveAutokeyConfigResponse{}
	obj.KeyProject = &project

	return obj, nil
}

func (r *autokeyAdminServer) GetKajPolicyConfigFolder(ctx context.Context, req *pb.GetKajPolicyConfigFolderRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKajPolicyConfigFolder not implemented")
}

func (r *autokeyAdminServer) UpdateKajPolicyConfigFolder(ctx context.Context, req *pb.UpdateKajPolicyConfigFolderRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateKajPolicyConfigFolder not implemented")
}

func (r *autokeyAdminServer) GetKajPolicyConfigOrganization(ctx context.Context, req *pb.GetKajPolicyConfigOrganizationRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKajPolicyConfigOrganization not implemented")
}

func (r *autokeyAdminServer) UpdateKajPolicyConfigOrganization(ctx context.Context, req *pb.UpdateKajPolicyConfigOrganizationRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateKajPolicyConfigOrganization not implemented")
}

func (r *autokeyAdminServer) GetKajPolicyConfigProject(ctx context.Context, req *pb.GetKajPolicyConfigProjectRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKajPolicyConfigProject not implemented")
}

func (r *autokeyAdminServer) UpdateKajPolicyConfigProject(ctx context.Context, req *pb.UpdateKajPolicyConfigProjectRequest) (*pb.KeyAccessJustificationsPolicyConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateKajPolicyConfigProject not implemented")
}

func (r *autokeyAdminServer) ShowEffectiveKeyAccessJustificationsPolicyConfigProject(ctx context.Context, req *pb.ShowEffectiveKeyAccessJustificationsPolicyConfigProjectRequest) (*pb.ShowEffectiveKeyAccessJustificationsPolicyConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowEffectiveKeyAccessJustificationsPolicyConfigProject not implemented")
}

type autokeyConfigName struct {
	folder string
}

func (a *autokeyConfigName) String() string {
	return "folders/" + a.folder + "/autokeyConfig"
}

// parseAutokeyConfigName parses a string into an AutoKeyConfig name.
// The expected form is `folders/{FOLDER_NUMBER}/autokeyConfig`.
func (r *autokeyAdminServer) parseAutokeyConfigName(name string) (*autokeyConfigName, error) {
	tokens := strings.Split(name, "/")
	if len(tokens) == 3 && tokens[0] == "folders" && tokens[2] == "autokeyConfig" {
		//fmt.Printf("Inside mock gcp controller %s\n\n", tokens[1])
		name := &autokeyConfigName{
			folder: tokens[1],
		}

		return name, nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "name %q is not valid", name)
}
