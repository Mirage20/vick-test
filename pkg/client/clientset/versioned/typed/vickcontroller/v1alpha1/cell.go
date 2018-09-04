/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/wso2/vick/pkg/apis/vickcontroller/v1alpha1"
	scheme "github.com/wso2/vick/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CellsGetter has a method to return a CellInterface.
// A group's client should implement this interface.
type CellsGetter interface {
	Cells(namespace string) CellInterface
}

// CellInterface has methods to work with Cell resources.
type CellInterface interface {
	Create(*v1alpha1.Cell) (*v1alpha1.Cell, error)
	Update(*v1alpha1.Cell) (*v1alpha1.Cell, error)
	UpdateStatus(*v1alpha1.Cell) (*v1alpha1.Cell, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Cell, error)
	List(opts v1.ListOptions) (*v1alpha1.CellList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Cell, err error)
	CellExpansion
}

// cells implements CellInterface
type cells struct {
	client rest.Interface
	ns     string
}

// newCells returns a Cells
func newCells(c *VickcontrollerV1alpha1Client, namespace string) *cells {
	return &cells{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cell, and returns the corresponding cell object, and an error if there is any.
func (c *cells) Get(name string, options v1.GetOptions) (result *v1alpha1.Cell, err error) {
	result = &v1alpha1.Cell{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cells").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Cells that match those selectors.
func (c *cells) List(opts v1.ListOptions) (result *v1alpha1.CellList, err error) {
	result = &v1alpha1.CellList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cells").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cells.
func (c *cells) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cells").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a cell and creates it.  Returns the server's representation of the cell, and an error, if there is any.
func (c *cells) Create(cell *v1alpha1.Cell) (result *v1alpha1.Cell, err error) {
	result = &v1alpha1.Cell{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cells").
		Body(cell).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cell and updates it. Returns the server's representation of the cell, and an error, if there is any.
func (c *cells) Update(cell *v1alpha1.Cell) (result *v1alpha1.Cell, err error) {
	result = &v1alpha1.Cell{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cells").
		Name(cell.Name).
		Body(cell).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cells) UpdateStatus(cell *v1alpha1.Cell) (result *v1alpha1.Cell, err error) {
	result = &v1alpha1.Cell{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cells").
		Name(cell.Name).
		SubResource("status").
		Body(cell).
		Do().
		Into(result)
	return
}

// Delete takes name of the cell and deletes it. Returns an error if one occurs.
func (c *cells) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cells").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cells) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cells").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cell.
func (c *cells) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Cell, err error) {
	result = &v1alpha1.Cell{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cells").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
