// Copyright 2020 Google LLC
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

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/apigateway/v1alpha1"
	scheme "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// APIGatewayAPIConfigsGetter has a method to return a APIGatewayAPIConfigInterface.
// A group's client should implement this interface.
type APIGatewayAPIConfigsGetter interface {
	APIGatewayAPIConfigs(namespace string) APIGatewayAPIConfigInterface
}

// APIGatewayAPIConfigInterface has methods to work with APIGatewayAPIConfig resources.
type APIGatewayAPIConfigInterface interface {
	Create(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.CreateOptions) (*v1alpha1.APIGatewayAPIConfig, error)
	Update(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.UpdateOptions) (*v1alpha1.APIGatewayAPIConfig, error)
	UpdateStatus(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.UpdateOptions) (*v1alpha1.APIGatewayAPIConfig, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.APIGatewayAPIConfig, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.APIGatewayAPIConfigList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.APIGatewayAPIConfig, err error)
	APIGatewayAPIConfigExpansion
}

// aPIGatewayAPIConfigs implements APIGatewayAPIConfigInterface
type aPIGatewayAPIConfigs struct {
	client rest.Interface
	ns     string
}

// newAPIGatewayAPIConfigs returns a APIGatewayAPIConfigs
func newAPIGatewayAPIConfigs(c *ApigatewayV1alpha1Client, namespace string) *aPIGatewayAPIConfigs {
	return &aPIGatewayAPIConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the aPIGatewayAPIConfig, and returns the corresponding aPIGatewayAPIConfig object, and an error if there is any.
func (c *aPIGatewayAPIConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.APIGatewayAPIConfig, err error) {
	result = &v1alpha1.APIGatewayAPIConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of APIGatewayAPIConfigs that match those selectors.
func (c *aPIGatewayAPIConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.APIGatewayAPIConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.APIGatewayAPIConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aPIGatewayAPIConfigs.
func (c *aPIGatewayAPIConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a aPIGatewayAPIConfig and creates it.  Returns the server's representation of the aPIGatewayAPIConfig, and an error, if there is any.
func (c *aPIGatewayAPIConfigs) Create(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.CreateOptions) (result *v1alpha1.APIGatewayAPIConfig, err error) {
	result = &v1alpha1.APIGatewayAPIConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aPIGatewayAPIConfig).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a aPIGatewayAPIConfig and updates it. Returns the server's representation of the aPIGatewayAPIConfig, and an error, if there is any.
func (c *aPIGatewayAPIConfigs) Update(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.UpdateOptions) (result *v1alpha1.APIGatewayAPIConfig, err error) {
	result = &v1alpha1.APIGatewayAPIConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		Name(aPIGatewayAPIConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aPIGatewayAPIConfig).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *aPIGatewayAPIConfigs) UpdateStatus(ctx context.Context, aPIGatewayAPIConfig *v1alpha1.APIGatewayAPIConfig, opts v1.UpdateOptions) (result *v1alpha1.APIGatewayAPIConfig, err error) {
	result = &v1alpha1.APIGatewayAPIConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		Name(aPIGatewayAPIConfig.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aPIGatewayAPIConfig).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the aPIGatewayAPIConfig and deletes it. Returns an error if one occurs.
func (c *aPIGatewayAPIConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aPIGatewayAPIConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched aPIGatewayAPIConfig.
func (c *aPIGatewayAPIConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.APIGatewayAPIConfig, err error) {
	result = &v1alpha1.APIGatewayAPIConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apigatewayapiconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}