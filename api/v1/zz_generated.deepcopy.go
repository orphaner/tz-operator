//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronStrategy) DeepCopyInto(out *CronStrategy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronStrategy.
func (in *CronStrategy) DeepCopy() *CronStrategy {
	if in == nil {
		return nil
	}
	out := new(CronStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CrossVersionObjectReference) DeepCopyInto(out *CrossVersionObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CrossVersionObjectReference.
func (in *CrossVersionObjectReference) DeepCopy() *CrossVersionObjectReference {
	if in == nil {
		return nil
	}
	out := new(CrossVersionObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlingResource) DeepCopyInto(out *IdlingResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlingResource.
func (in *IdlingResource) DeepCopy() *IdlingResource {
	if in == nil {
		return nil
	}
	out := new(IdlingResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IdlingResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlingResourceList) DeepCopyInto(out *IdlingResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IdlingResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlingResourceList.
func (in *IdlingResourceList) DeepCopy() *IdlingResourceList {
	if in == nil {
		return nil
	}
	out := new(IdlingResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IdlingResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlingResourceSpec) DeepCopyInto(out *IdlingResourceSpec) {
	*out = *in
	out.IdlingResourceRef = in.IdlingResourceRef
	if in.ResumeReplicas != nil {
		in, out := &in.ResumeReplicas, &out.ResumeReplicas
		*out = new(int32)
		**out = **in
	}
	if in.IdlingStrategy != nil {
		in, out := &in.IdlingStrategy, &out.IdlingStrategy
		*out = new(IdlingStrategy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlingResourceSpec.
func (in *IdlingResourceSpec) DeepCopy() *IdlingResourceSpec {
	if in == nil {
		return nil
	}
	out := new(IdlingResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlingResourceStatus) DeepCopyInto(out *IdlingResourceStatus) {
	*out = *in
	if in.PreviousReplicas != nil {
		in, out := &in.PreviousReplicas, &out.PreviousReplicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlingResourceStatus.
func (in *IdlingResourceStatus) DeepCopy() *IdlingResourceStatus {
	if in == nil {
		return nil
	}
	out := new(IdlingResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IdlingStrategy) DeepCopyInto(out *IdlingStrategy) {
	*out = *in
	if in.CronStrategy != nil {
		in, out := &in.CronStrategy, &out.CronStrategy
		*out = new(CronStrategy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IdlingStrategy.
func (in *IdlingStrategy) DeepCopy() *IdlingStrategy {
	if in == nil {
		return nil
	}
	out := new(IdlingStrategy)
	in.DeepCopyInto(out)
	return out
}
