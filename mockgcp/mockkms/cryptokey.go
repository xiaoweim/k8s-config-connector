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
// proto.service: google.cloud.kms.v1.KeyManagementService
// proto.message: google.cloud.kms.v1.CryptoKey

package mockkms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/kms/v1"
)

func (r *kmsServer) GetProjectsLocationsKeyRingsCryptoKey(ctx context.Context, req *pb.GetProjectsLocationsKeyRingsCryptoKeyRequest) (*pb.CryptoKey, error) {
	name, err := r.parseCryptoKeyName(req.GetName())
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	obj := &pb.CryptoKey{}
	if err := r.storage.Get(ctx, fqn, obj); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "CryptoKey %s not found.", fqn)
		}
		return nil, err
	}

	return obj, nil
}

func (r *kmsServer) CreateProjectsLocationsKeyRingsCryptoKey(ctx context.Context, req *pb.CreateProjectsLocationsKeyRingsCryptoKeyRequest) (*pb.CryptoKey, error) {
	reqName := fmt.Sprintf("%s/cryptoKeys/%s", req.GetParent(), req.GetCryptoKeyId())
	name, err := r.parseCryptoKeyName(reqName)
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	now := time.Now()

	obj := proto.Clone(req.GetProjectsLocationsKeyRingsCryptoKey()).(*pb.CryptoKey)
	obj.Name = &fqn
	obj.CreateTime = timestamppb.New(now)

	r.populateDefaultsForCryptoKey(name, obj)

	if !req.GetSkipInitialVersionCreation() {
		var primary *pb.CryptoKeyVersion

		if obj.VersionTemplate != nil {
			primary = &pb.CryptoKeyVersion{
				Algorithm:       pb.CryptoKeyVersion_CryptoKeyVersionAlgorithm(int32(obj.VersionTemplate.GetAlgorithm())).Enum(),
				ProtectionLevel: pb.CryptoKeyVersion_ProtectionLevel(int32(obj.VersionTemplate.GetProtectionLevel())).Enum(),
			}
		} else if req.GetProjectsLocationsKeyRingsCryptoKey().GetPurpose() == pb.CryptoKey_ENCRYPT_DECRYPT {
			// Set default
			primary = &pb.CryptoKeyVersion{
				Algorithm:       pb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION.Enum(),
				ProtectionLevel: pb.CryptoKeyVersion_SOFTWARE.Enum(),
			}
		} else {

			primary = &pb.CryptoKeyVersion{
				// Algorithm is required
				Algorithm: pb.CryptoKeyVersion_CryptoKeyVersionAlgorithm(int32(obj.VersionTemplate.GetAlgorithm())).Enum(),
			}
		}
		createVersionReq := &pb.CreateProjectsLocationsKeyRingsCryptoKeysCryptoKeyVersionRequest{
			Parent: &fqn,
			ProjectsLocationsKeyRingsCryptoKeysCryptoKeyVersion: primary,
		}
		createdVersion, err := r.CreateProjectsLocationsKeyRingsCryptoKeysCryptoKeyVersion(ctx, createVersionReq)
		if err != nil {
			return nil, err
		}
		obj.Primary = createdVersion
		obj.VersionTemplate = &pb.CryptoKeyVersionTemplate{
			Algorithm:       pb.CryptoKeyVersion_CryptoKeyVersionAlgorithm(int32(createdVersion.GetAlgorithm())).Enum(),
			ProtectionLevel: pb.CryptoKeyVersion_ProtectionLevel(int32(createdVersion.GetProtectionLevel())).Enum(),
		}
	}
	if err := r.storage.Create(ctx, fqn, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *kmsServer) PatchProjectsLocationsKeyRingsCryptoKey(ctx context.Context, req *pb.PatchProjectsLocationsKeyRingsCryptoKeyRequest) (*pb.CryptoKey, error) {
	name, err := r.parseCryptoKeyName(req.GetProjectsLocationsKeyRingsCryptoKey().GetName())
	if err != nil {
		return nil, err
	}

	fqn := name.String()

	obj := &pb.CryptoKey{}
	if err := r.storage.Get(ctx, fqn, obj); err != nil {
		return nil, err
	}

	updateMask := req.GetUpdateMask()
	paths := strings.Split(updateMask, ",")
	for _, path := range paths {
		switch path {
		case "labels":
			obj.Labels = req.GetProjectsLocationsKeyRingsCryptoKey().GetLabels()
		default:
			return nil, status.Errorf(codes.InvalidArgument, "field %q is not yet handled in mock", path)
		}
	}

	if err := r.storage.Update(ctx, fqn, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *kmsServer) populateDefaultsForCryptoKey(name *CryptoKeyName, obj *pb.CryptoKey) {
	if obj.DestroyScheduledDuration == nil {
		obj.DestroyScheduledDuration = durationpb.New(time.Hour * 24 * 30)
	}
}

type CryptoKeyName struct {
	KeyRingName
	CryptoKeyID string
}

func (n *CryptoKeyName) String() string {
	return n.KeyRingName.String() + "/cryptoKeys/" + n.CryptoKeyID
}

// parseCryptoKeyName parses a string into an CryptoKeyName.
// The expected form is `projects/*/locations/*/keyRings/*/cryptoKeys/*`.
func (r *kmsServer) parseCryptoKeyName(name string) (*CryptoKeyName, error) {
	tokens := strings.Split(name, "/")

	if len(tokens) == 8 && tokens[6] == "cryptoKeys" {
		keyRingName, err := r.parseKeyRingName(strings.Join(tokens[0:6], "/"))
		if err != nil {
			return nil, err
		}

		name := &CryptoKeyName{
			KeyRingName: *keyRingName,
			CryptoKeyID: tokens[7],
		}

		return name, nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "name %q is not valid", name)
}
