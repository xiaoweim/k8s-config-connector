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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/controllerbuilder/pkg/codegen"
	"github.com/GoogleCloudPlatform/k8s-config-connector/dev/tools/controllerbuilder/pkg/gocode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// List of proto message types not mapped to Go structs
var protoMsgsNotMappedToGoStruct = []string{
	"google.protobuf.Timestamp",
	"google.protobuf.Duration",
	"google.protobuf.Int64Value",
}

func main() {
	config := parseFlags()
	updator := NewTypeUpdater(config, protoMsgsNotMappedToGoStruct)

	if err := updator.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type TypeUpdater struct {
	opts                         options
	protoMsgsNotMappedToGoStruct []string
}

type options struct {
	protoSourcePath       string
	parentMessageFullName string
	newField              string
	ignoredFields         string
	apiDirectory          string
	goPackagePath         string
}

func parseFlags() options {
	var opts options
	flag.StringVar(&opts.protoSourcePath, "proto-source-path", "../../../proto-to-mapper/build/googleapis.pb", "Path to (compiled) proto for APIs")
	flag.StringVar(&opts.parentMessageFullName, "parent-message-full-name", opts.parentMessageFullName, "Fully qualified name of the proto message holding the new field")
	flag.StringVar(&opts.newField, "new-field", opts.newField, "Name of the new field")
	flag.StringVar(&opts.goPackagePath, "api-go-package-path", "github.com/GoogleCloudPlatform/k8s-config-connector/apis/", "Package path")
	flag.StringVar(&opts.apiDirectory, "api-dir", "../../../../../apis/", "Base directory for APIs")
	flag.StringVar(&opts.ignoredFields, "ignored-fields", opts.ignoredFields, "Comma-separated list of fields to ignore")
	flag.Parse()
	return opts
}

func NewTypeUpdater(opts options, protoMsgsNotMappedToGoStruct []string) *TypeUpdater {
	return &TypeUpdater{
		opts:                         opts,
		protoMsgsNotMappedToGoStruct: protoMsgsNotMappedToGoStruct,
	}
}

func (u *TypeUpdater) Run() error {
	parentMsg, newField, err := findNewField(u.opts.protoSourcePath, u.opts.parentMessageFullName, u.opts.newField)
	if err != nil {
		return err
	}

	msgs, err := findDependentMsgs(newField, strings.Split(u.opts.ignoredFields, ","))
	if err != nil {
		return err
	}

	removeNotMappedToGoStruct(msgs)

	if err := removeAlreadyGenerated(u.opts.goPackagePath, u.opts.apiDirectory, msgs); err != nil {
		return err
	}

	// TODO: find the proper file to write/update. Currently print to a tmp file.
	return writeToLocalFile(newField, parentMsg, msgs)
}

// findNewField locates the parent message and the new field in the proto file
func findNewField(pbSourcePath, parentMsgFullName, newFieldName string) (protoreflect.MessageDescriptor, protoreflect.FieldDescriptor, error) {
	fileData, err := os.ReadFile(pbSourcePath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %q: %w", pbSourcePath, err)
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(fileData, fds); err != nil {
		return nil, nil, fmt.Errorf("unmarshalling %q: %w", pbSourcePath, err)
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, nil, fmt.Errorf("building file description: %w", err)
	}

	// Find the parent message
	messageDesc, err := files.FindDescriptorByName(protoreflect.FullName(parentMsgFullName))
	if err != nil {
		return nil, nil, err
	}
	msgDesc, ok := messageDesc.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected descriptor type: %T", msgDesc)
	}

	// Find the new field in parent message
	fieldDesc := msgDesc.Fields().ByName(protoreflect.Name(newFieldName))
	if fieldDesc == nil {
		return nil, nil, fmt.Errorf("field not found in message")
	}

	return msgDesc, fieldDesc, nil
}

// findDependentMsgs finds all dependent proto messages for the given field, ignoring specified fields
func findDependentMsgs(field protoreflect.FieldDescriptor, ignoredProtoFields []string) (map[string]protoreflect.MessageDescriptor, error) {
	deps := make(map[string]protoreflect.MessageDescriptor)
	analyzeDependencies(field, deps, ignoredProtoFields)

	// fmt.Println("\nDependencies:")
	// for _, dep := range deps {
	// 	fmt.Printf("%s\n", dep.FullName())
	// }
	return deps, nil
}

// analyzeDependencies recursively explores the dependent proto messages of the given field.
func analyzeDependencies(field protoreflect.FieldDescriptor, deps map[string]protoreflect.MessageDescriptor, ignoreFields []string) {
	if shouldIgnore(string(field.FullName()), ignoreFields) {
		return
	}

	if field.IsMap() {
		/* an example Map field
		message ParentMessage {
			map<KeyMessage, ValueMessage> some_map_field = 1;
		} */
		mapEntry := field.Message()
		if keyField := mapEntry.Fields().ByName("key"); keyField != nil {
			analyzeDependencies(keyField, deps, ignoreFields)
		}
		if valueField := mapEntry.Fields().ByName("value"); valueField != nil {
			analyzeDependencies(valueField, deps, ignoreFields)
		}
	} else {
		switch field.Kind() {
		case protoreflect.MessageKind:
			msg := field.Message()
			fqn := string(msg.FullName())
			if _, ok := deps[fqn]; !ok {
				deps[fqn] = msg
				for i := 0; i < msg.Fields().Len(); i++ {
					field := msg.Fields().Get(i)
					analyzeDependencies(field, deps, ignoreFields)
				}
			}
		case protoreflect.EnumKind:
			// deps[string(field.Enum().FullName())] = true  // Skip enum because enum is mapped to Go string in code generation
		}
	}
}

func shouldIgnore(fieldName string, ignoreFields []string) bool {
	for _, ignoreField := range ignoreFields {
		if fieldName == ignoreField {
			return true
		}
	}
	return false
}

func removeNotMappedToGoStruct(msgs map[string]protoreflect.MessageDescriptor) {
	for _, msgName := range protoMsgsNotMappedToGoStruct {
		delete(msgs, msgName)
	}
}

// removeAlreadyGenerated removes proto messages that have already been generated (including manually edited)
func removeAlreadyGenerated(goPackagePath, outputAPIDirectory string, targets map[string]protoreflect.MessageDescriptor) error {
	packages, err := gocode.LoadPackageTree(goPackagePath, outputAPIDirectory)
	if err != nil {
		return err
	}
	for _, p := range packages {
		for _, s := range p.Structs {
			if annotation := s.GetAnnotation("+kcc:proto"); annotation != "" {
				if _, ok := targets[annotation]; ok {
					delete(targets, annotation)
				}
			}
		}
	}
	return nil
}

func writeToLocalFile(field protoreflect.FieldDescriptor, parentMsg protoreflect.MessageDescriptor, msgs map[string]protoreflect.MessageDescriptor) error {
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "[The generated new field]\n\n")
	codegen.WriteField(file, field, parentMsg, 0)

	fmt.Fprintf(file, "\n\n[The generated new messages]\n\n")
	for _, msg := range msgs {
		codegen.WriteMessage(file, msg)
	}
	return nil
}