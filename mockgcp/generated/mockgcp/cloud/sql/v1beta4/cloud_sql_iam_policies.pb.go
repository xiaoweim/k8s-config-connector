// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: mockgcp/cloud/sql/v1beta4/cloud_sql_iam_policies.proto

package sqlpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto protoreflect.FileDescriptor

var file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_rawDesc = []byte{
	0x0a, 0x36, 0x6d, 0x6f, 0x63, 0x6b, 0x67, 0x63, 0x70, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f,
	0x73, 0x71, 0x6c, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x34, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x5f, 0x73, 0x71, 0x6c, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x6d, 0x6f, 0x63, 0x6b, 0x67, 0x63,
	0x70, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x73, 0x71, 0x6c, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x34, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x33, 0x0a, 0x15, 0x53, 0x71,
	0x6c, 0x49, 0x61, 0x6d, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x1a, 0x1a, 0xca, 0x41, 0x17, 0x73, 0x71, 0x6c, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x42,
	0x6b, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x6f, 0x63, 0x6b, 0x67, 0x63, 0x70, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x73, 0x71, 0x6c, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x34,
	0x42, 0x18, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x53, 0x71, 0x6c, 0x49, 0x61, 0x6d, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x69, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x6f, 0x2f, 0x73, 0x71, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x34,
	0x2f, 0x73, 0x71, 0x6c, 0x70, 0x62, 0x3b, 0x73, 0x71, 0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_goTypes = []interface{}{}
var file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_init() }
func file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_init() {
	if File_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_goTypes,
		DependencyIndexes: file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_depIdxs,
	}.Build()
	File_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto = out.File
	file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_rawDesc = nil
	file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_goTypes = nil
	file_mockgcp_cloud_sql_v1beta4_cloud_sql_iam_policies_proto_depIdxs = nil
}
