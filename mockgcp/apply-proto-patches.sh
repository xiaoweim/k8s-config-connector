#!/bin/bash

# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -o errexit
set -o nounset
set -o pipefail

set -x

REPO_ROOT=$(git rev-parse --show-toplevel)
cd ${REPO_ROOT}/mockgcp

cd tools/patch-proto

go run . --file ${REPO_ROOT}/mockgcp/apis/mockgcp/cloud/apigee/v1/service.proto --service "ProjectsServer" --mode "replace" <<EOF
  // Provisions a new Apigee organization with a functioning runtime. This is the standard way to create trial organizations for a free Apigee trial.
  rpc ProvisionOrganizationProject(ProvisionOrganizationProjectRequest) returns (.google.longrunning.Operation) {
    option (google.api.http) = {
      post: "/v1/{name=projects/*}:provisionOrganization"
      body: "project"
    };
  };
EOF

# SQL patches

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/sql/v1beta4/cloud_sql_resources.proto --message BackupConfiguration --mode append <<EOF

  // The tier of the backup.

  optional string backup_tier = 14;

EOF

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/sql/v1beta4/cloud_sql_resources.proto --message IpConfiguration --mode append <<EOF

  // The server certificate rotation mode.

  optional string server_certificate_rotation_mode = 10;

EOF

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/sql/v1beta4/cloud_sql_resources.proto --message Settings --mode append <<EOF

  // Maximum replication lag in seconds.

  optional int32 replication_lag_max_seconds = 100;

EOF

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/sql/v1beta4/cloud_sql_resources.proto --message DatabaseInstance --mode append <<EOF

  // Whether to include replicas for major version upgrade.

  optional bool include_replicas_for_major_version_upgrade = 56;



  // Whether the instance satisfies PZI.

  optional bool satisfies_pzi = 57;

EOF

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/sql/v1beta4/cloud_sql_users.proto --message User --mode append <<EOF

  // The status of the user's IAM authentication.

  optional string iam_status = 16;

EOF



# AlloyDB patches

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/cloud/alloydb/v1beta/resources.proto --message ConnectionPoolConfig --mode append <<EOF

  // The number of pooler instances.

  int32 pooler_count = 14;

EOF

# Container/GKE patches

go run . --file ${REPO_ROOT}/mockgcp/third_party/googleapis/google/container/v1beta1/cluster_service.proto --message DNSEndpointConfig --mode append <<EOF

  // Controls whether the k8s token auth is allowed via DNS.

  optional bool enable_k8s_tokens_via_dns = 5;

EOF

# ResourceManager v1 patches - temporarily switching to proto3 because patch-proto has issues with proto2
API_PROTO=${REPO_ROOT}/mockgcp/apis/mockgcp/cloud/resourcemanager/v1/api.proto
sed -i 's/^syntax = "proto2";/syntax = "proto3";/' ${API_PROTO}

go run . --file ${API_PROTO} --message "TestIamPermissionsRequest" --mode "replace" <<EOF
  // REQUIRED: The resource for which the permission checking is to be performed.
  optional string resource = 1;

  // The set of permissions to check for the resource.
  repeated string permissions = 2 [json_name="permissions"];
EOF

go run . --file ${API_PROTO} --service "FoldersServer" --mode "append" <<EOF
  // Returns permissions that a caller has on the specified project.
  rpc TestIamPermissions(TestIamPermissionsRequest) returns (TestIamPermissionsResponse) {
    option (google.api.http) = {
      post: "/v1/{resource=folders/*}:testIamPermissions"
      body: "*"
    };
  };
EOF

go run . --file ${API_PROTO} --service "OrganizationsServer" --mode "append" <<EOF
  // Returns permissions that a caller has on the specified project.
  rpc TestIamPermissions(TestIamPermissionsRequest) returns (TestIamPermissionsResponse) {
    option (google.api.http) = {
      post: "/v1/{resource=organizations/*}:testIamPermissions"
      body: "*"
    };
  };
EOF

go run . --file ${API_PROTO} --service "ProjectsServer" --mode "append" <<EOF
  // Returns permissions that a caller has on the specified project.
  rpc TestIamPermissions(TestIamPermissionsRequest) returns (TestIamPermissionsResponse) {
    option (google.api.http) = {
      post: "/v1/{resource=projects/*}:testIamPermissions"
      body: "*"
    };
  };
EOF

sed -i 's/^syntax = "proto3";/syntax = "proto2";/' ${API_PROTO}

# KMS patches - temporarily switching to proto3 because patch-proto has issues with proto2
API_PROTO=${REPO_ROOT}/mockgcp/apis/mockgcp/cloud/kms/v1/api.proto
sed -i 's/^syntax = "proto2";/syntax = "proto3";/' ${API_PROTO}

go run . --file ${API_PROTO} --message "ShowEffectiveKeyAccessJustificationsEnrollmentConfigProjectRequest" --mode "replace" <<EOF
  optional string project = 1;
EOF

go run . --file ${API_PROTO} --message "ShowEffectiveKeyAccessJustificationsPolicyConfigProjectRequest" --mode "replace" <<EOF
  optional string project = 1;
EOF

go run . --file ${API_PROTO} --message "AutokeyConfig" --mode "replace" <<EOF
  enum AutokeyConfigState {
    AUTOKEY_CONFIG_STATE_UNSPECIFIED = 0;
    ACTIVE = 1;
    KEY_PROJECT_DELETED = 2;
    UNINITIALIZED = 3;
  }
  optional string etag = 1 [json_name="etag"];
  optional string key_project = 2 [json_name="keyProject"];
  optional string key_project_resolution_mode = 3 [json_name="keyProjectResolutionMode"];
  optional string name = 4 [json_name="name"];
  optional AutokeyConfigState state = 5 [json_name="state"];
EOF

go run . --file ${API_PROTO} --message "KeyRing" --mode "replace" <<EOF
  optional string name = 1 [json_name="name"];
  optional .google.protobuf.Timestamp create_time = 2 [json_name="createTime"];
EOF

go run . --file ${API_PROTO} --message "KeyHandle" --mode "replace" <<EOF
  optional string name = 1 [json_name="name"];
  optional string kms_key = 3 [json_name="kmsKey"];
  optional string resource_type_selector = 4 [json_name="resourceTypeSelector"];
EOF

cat >> ${API_PROTO} <<EOF
message CreateKeyHandleMetadata {
}
EOF

# Patch enums for KMS with CORRECT field numbers from official protos
go run . --file ${API_PROTO} --message "CryptoKey" --mode "replace" <<EOF
  enum CryptoKeyPurpose {
    CRYPTO_KEY_PURPOSE_UNSPECIFIED = 0;
    ENCRYPT_DECRYPT = 1;
    ASYMMETRIC_SIGN = 5;
    ASYMMETRIC_DECRYPT = 6;
    RAW_ENCRYPT_DECRYPT = 7;
    MAC = 9;
  }
  optional string name = 1 [json_name="name"];
  optional CryptoKeyVersion primary = 2 [json_name="primary"];
  optional CryptoKeyPurpose purpose = 3 [json_name="purpose"];
  optional .google.protobuf.Timestamp create_time = 5 [json_name="createTime"];
  optional .google.protobuf.Timestamp next_rotation_time = 7 [json_name="nextRotationTime"];
  optional .google.protobuf.Duration rotation_period = 8 [json_name="rotationPeriod"];
  optional CryptoKeyVersionTemplate version_template = 11 [json_name="versionTemplate"];
  map<string, string> labels = 12 [json_name="labels"];
  optional bool import_only = 13 [json_name="importOnly"];
  optional .google.protobuf.Duration destroy_scheduled_duration = 14 [json_name="destroyScheduledDuration"];
  optional string crypto_key_backend = 15 [json_name="cryptoKeyBackend"];
  optional KeyAccessJustificationsPolicy key_access_justifications_policy = 17 [json_name="keyAccessJustificationsPolicy"];
EOF

go run . --file ${API_PROTO} --message "CryptoKeyVersion" --mode "replace" <<EOF
  enum CryptoKeyVersionAlgorithm {
    CRYPTO_KEY_VERSION_ALGORITHM_UNSPECIFIED = 0;
    GOOGLE_SYMMETRIC_ENCRYPTION = 1;
    RSA_SIGN_PSS_2048_SHA256 = 2;
    RSA_SIGN_PSS_3072_SHA256 = 3;
    RSA_SIGN_PSS_4096_SHA256 = 4;
    RSA_SIGN_PSS_4096_SHA512 = 15;
    RSA_SIGN_PKCS1_2048_SHA256 = 5;
    RSA_SIGN_PKCS1_3072_SHA256 = 6;
    RSA_SIGN_PKCS1_4096_SHA256 = 7;
    RSA_SIGN_PKCS1_4096_SHA512 = 16;
    RSA_SIGN_RAW_PKCS1_2048 = 28;
    RSA_SIGN_RAW_PKCS1_3072 = 29;
    RSA_SIGN_RAW_PKCS1_4096 = 30;
    RSA_DECRYPT_OAEP_2048_SHA256 = 8;
    RSA_DECRYPT_OAEP_3072_SHA256 = 9;
    RSA_DECRYPT_OAEP_4096_SHA256 = 10;
    RSA_DECRYPT_OAEP_4096_SHA512 = 17;
    RSA_DECRYPT_OAEP_2048_SHA1 = 37;
    RSA_DECRYPT_OAEP_3072_SHA1 = 38;
    RSA_DECRYPT_OAEP_4096_SHA1 = 39;
    EC_SIGN_P256_SHA256 = 12;
    EC_SIGN_P384_SHA384 = 13;
    EC_SIGN_SECP256K1_SHA256 = 31;
    HMAC_SHA256 = 32;
    HMAC_SHA1 = 33;
    HMAC_SHA384 = 34;
    HMAC_SHA512 = 35;
    HMAC_SHA224 = 36;
    EXTERNAL_SYMMETRIC_ENCRYPTION = 18;
  }
  enum CryptoKeyVersionState {
    CRYPTO_KEY_VERSION_STATE_UNSPECIFIED = 0;
    ENABLED = 1;
    DISABLED = 2;
    DESTROYED = 3;
    DESTROY_SCHEDULED = 4;
    PENDING_GENERATION = 5;
    PENDING_IMPORT = 6;
    IMPORT_FAILED = 7;
    GENERATION_FAILED = 8;
    PENDING_EXTERNAL_DESTRUCTION = 9;
    EXTERNAL_DESTRUCTION_FAILED = 10;
  }
  enum ProtectionLevel {
    PROTECTION_LEVEL_UNSPECIFIED = 0;
    SOFTWARE = 1;
    HSM = 2;
    EXTERNAL = 3;
    EXTERNAL_VPC = 4;
  }
  optional string name = 1 [json_name="name"];
  optional CryptoKeyVersionState state = 3 [json_name="state"];
  optional .google.protobuf.Timestamp create_time = 4 [json_name="createTime"];
  optional .google.protobuf.Timestamp destroy_time = 5 [json_name="destroyTime"];
  optional .google.protobuf.Timestamp destroy_event_time = 6 [json_name="destroyEventTime"];
  optional ProtectionLevel protection_level = 7 [json_name="protectionLevel"];
  optional KeyOperationAttestation attestation = 8 [json_name="attestation"];
  optional CryptoKeyVersionAlgorithm algorithm = 10 [json_name="algorithm"];
  optional .google.protobuf.Timestamp generate_time = 11 [json_name="generateTime"];
  optional string import_job = 14 [json_name="importJob"];
  optional .google.protobuf.Timestamp import_time = 15 [json_name="importTime"];
  optional string import_failure_reason = 16 [json_name="importFailureReason"];
  optional string generation_failure_reason = 19 [json_name="generationFailureReason"];
  optional .google.protobuf.Timestamp external_destruction_time = 20 [json_name="externalDestructionTime"];
  optional .google.protobuf.Timestamp external_destruction_event_time = 21 [json_name="externalDestructionEventTime"];
  optional bool reimport_eligible = 24 [json_name="reimportEligible"];
EOF

go run . --file ${API_PROTO} --message "CryptoKeyVersionTemplate" --mode "replace" <<EOF
  optional .mockgcp.cloud.kms.v1.CryptoKeyVersion.CryptoKeyVersionAlgorithm algorithm = 1 [json_name="algorithm"];
  optional .mockgcp.cloud.kms.v1.CryptoKeyVersion.ProtectionLevel protection_level = 2 [json_name="protectionLevel"];
EOF

go run . --file ${API_PROTO} --message "ImportJob" --mode "replace" <<EOF
  enum ImportMethod {
    IMPORT_METHOD_UNSPECIFIED = 0;
    RSA_OAEP_3072_SHA1_AES_256 = 1;
    RSA_OAEP_4096_SHA1_AES_256 = 2;
    RSA_OAEP_3072_SHA256_AES_256 = 3;
    RSA_OAEP_4096_SHA256_AES_256 = 4;
    RSA_OAEP_3072_SHA256 = 5;
    RSA_OAEP_4096_SHA256 = 6;
  }
  enum ImportJobState {
    IMPORT_JOB_STATE_UNSPECIFIED = 0;
    PENDING_GENERATION = 1;
    ACTIVE = 2;
    EXPIRED = 3;
  }
  optional string name = 1 [json_name="name"];
  optional ImportMethod import_method = 2 [json_name="importMethod"];
  optional .mockgcp.cloud.kms.v1.CryptoKeyVersion.ProtectionLevel protection_level = 3 [json_name="protectionLevel"];
  optional .google.protobuf.Timestamp create_time = 4 [json_name="createTime"];
  optional .google.protobuf.Timestamp generate_time = 5 [json_name="generateTime"];
  optional ImportJobState state = 6 [json_name="state"];
  optional .google.protobuf.Timestamp expire_time = 7 [json_name="expireTime"];
  optional .google.protobuf.Timestamp expire_event_time = 8 [json_name="expireEventTime"];
  optional KeyOperationAttestation attestation = 9 [json_name="attestation"];
  optional WrappingPublicKey public_key = 10 [json_name="publicKey"];
  optional string crypto_key_backend = 11 [json_name="cryptoKeyBackend"];
EOF

go run . --file ${API_PROTO} --message "KeyOperationAttestation" --mode "replace" <<EOF
  enum AttestationFormat {
    ATTESTATION_FORMAT_UNSPECIFIED = 0;
    CAVIUM_V1_COMPRESSED = 3;
    CAVIUM_V2_COMPRESSED = 4;
  }
  optional AttestationFormat format = 4 [json_name="format"];
  optional bytes content = 5 [json_name="content"];
  optional CertificateChains cert_chains = 6 [json_name="certChains"];
EOF

sed -i 's/^syntax = "proto3";/syntax = "proto2";/' ${API_PROTO}
