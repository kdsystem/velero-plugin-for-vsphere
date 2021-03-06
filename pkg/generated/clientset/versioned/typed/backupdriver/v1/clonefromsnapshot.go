/*
Copyright the Velero contributors.

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

package v1

import (
	"time"

	v1 "github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/apis/backupdriver/v1"
	scheme "github.com/vmware-tanzu/velero-plugin-for-vsphere/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CloneFromSnapshotsGetter has a method to return a CloneFromSnapshotInterface.
// A group's client should implement this interface.
type CloneFromSnapshotsGetter interface {
	CloneFromSnapshots(namespace string) CloneFromSnapshotInterface
}

// CloneFromSnapshotInterface has methods to work with CloneFromSnapshot resources.
type CloneFromSnapshotInterface interface {
	Create(*v1.CloneFromSnapshot) (*v1.CloneFromSnapshot, error)
	Update(*v1.CloneFromSnapshot) (*v1.CloneFromSnapshot, error)
	UpdateStatus(*v1.CloneFromSnapshot) (*v1.CloneFromSnapshot, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.CloneFromSnapshot, error)
	List(opts metav1.ListOptions) (*v1.CloneFromSnapshotList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloneFromSnapshot, err error)
	CloneFromSnapshotExpansion
}

// cloneFromSnapshots implements CloneFromSnapshotInterface
type cloneFromSnapshots struct {
	client rest.Interface
	ns     string
}

// newCloneFromSnapshots returns a CloneFromSnapshots
func newCloneFromSnapshots(c *BackupdriverV1Client, namespace string) *cloneFromSnapshots {
	return &cloneFromSnapshots{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cloneFromSnapshot, and returns the corresponding cloneFromSnapshot object, and an error if there is any.
func (c *cloneFromSnapshots) Get(name string, options metav1.GetOptions) (result *v1.CloneFromSnapshot, err error) {
	result = &v1.CloneFromSnapshot{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CloneFromSnapshots that match those selectors.
func (c *cloneFromSnapshots) List(opts metav1.ListOptions) (result *v1.CloneFromSnapshotList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CloneFromSnapshotList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cloneFromSnapshots.
func (c *cloneFromSnapshots) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a cloneFromSnapshot and creates it.  Returns the server's representation of the cloneFromSnapshot, and an error, if there is any.
func (c *cloneFromSnapshots) Create(cloneFromSnapshot *v1.CloneFromSnapshot) (result *v1.CloneFromSnapshot, err error) {
	result = &v1.CloneFromSnapshot{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		Body(cloneFromSnapshot).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cloneFromSnapshot and updates it. Returns the server's representation of the cloneFromSnapshot, and an error, if there is any.
func (c *cloneFromSnapshots) Update(cloneFromSnapshot *v1.CloneFromSnapshot) (result *v1.CloneFromSnapshot, err error) {
	result = &v1.CloneFromSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		Name(cloneFromSnapshot.Name).
		Body(cloneFromSnapshot).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cloneFromSnapshots) UpdateStatus(cloneFromSnapshot *v1.CloneFromSnapshot) (result *v1.CloneFromSnapshot, err error) {
	result = &v1.CloneFromSnapshot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		Name(cloneFromSnapshot.Name).
		SubResource("status").
		Body(cloneFromSnapshot).
		Do().
		Into(result)
	return
}

// Delete takes name of the cloneFromSnapshot and deletes it. Returns an error if one occurs.
func (c *cloneFromSnapshots) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cloneFromSnapshots) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cloneFromSnapshot.
func (c *cloneFromSnapshots) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloneFromSnapshot, err error) {
	result = &v1.CloneFromSnapshot{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("clonefromsnapshots").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
