//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomerManagedEncryption) DeepCopyInto(out *CustomerManagedEncryption) {
	*out = *in
	out.KmsKeyRef = in.KmsKeyRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomerManagedEncryption.
func (in *CustomerManagedEncryption) DeepCopy() *CustomerManagedEncryption {
	if in == nil {
		return nil
	}
	out := new(CustomerManagedEncryption)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replication) DeepCopyInto(out *Replication) {
	*out = *in
	if in.LegacyAuto != nil {
		in, out := &in.LegacyAuto, &out.LegacyAuto
		*out = new(bool)
		**out = **in
	}
	if in.LegacyAutomatic != nil {
		in, out := &in.LegacyAutomatic, &out.LegacyAutomatic
		*out = new(Replication_Automatic)
		(*in).DeepCopyInto(*out)
	}
	if in.UserManaged != nil {
		in, out := &in.UserManaged, &out.UserManaged
		*out = new(Replication_UserManaged)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replication.
func (in *Replication) DeepCopy() *Replication {
	if in == nil {
		return nil
	}
	out := new(Replication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replication_Automatic) DeepCopyInto(out *Replication_Automatic) {
	*out = *in
	if in.CustomerManagedEncryption != nil {
		in, out := &in.CustomerManagedEncryption, &out.CustomerManagedEncryption
		*out = new(CustomerManagedEncryption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replication_Automatic.
func (in *Replication_Automatic) DeepCopy() *Replication_Automatic {
	if in == nil {
		return nil
	}
	out := new(Replication_Automatic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replication_UserManaged) DeepCopyInto(out *Replication_UserManaged) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = make([]Replication_UserManaged_Replica, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replication_UserManaged.
func (in *Replication_UserManaged) DeepCopy() *Replication_UserManaged {
	if in == nil {
		return nil
	}
	out := new(Replication_UserManaged)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replication_UserManaged_Replica) DeepCopyInto(out *Replication_UserManaged_Replica) {
	*out = *in
	if in.Location != nil {
		in, out := &in.Location, &out.Location
		*out = new(string)
		**out = **in
	}
	if in.CustomerManagedEncryption != nil {
		in, out := &in.CustomerManagedEncryption, &out.CustomerManagedEncryption
		*out = new(CustomerManagedEncryption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replication_UserManaged_Replica.
func (in *Replication_UserManaged_Replica) DeepCopy() *Replication_UserManaged_Replica {
	if in == nil {
		return nil
	}
	out := new(Replication_UserManaged_Replica)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rotation) DeepCopyInto(out *Rotation) {
	*out = *in
	if in.NextRotationTime != nil {
		in, out := &in.NextRotationTime, &out.NextRotationTime
		*out = new(string)
		**out = **in
	}
	if in.RotationPeriod != nil {
		in, out := &in.RotationPeriod, &out.RotationPeriod
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rotation.
func (in *Rotation) DeepCopy() *Rotation {
	if in == nil {
		return nil
	}
	out := new(Rotation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secret) DeepCopyInto(out *Secret) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Replication != nil {
		in, out := &in.Replication, &out.Replication
		*out = new(Replication)
		(*in).DeepCopyInto(*out)
	}
	if in.CreateTime != nil {
		in, out := &in.CreateTime, &out.CreateTime
		*out = new(string)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Topics != nil {
		in, out := &in.Topics, &out.Topics
		*out = make([]Topic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ExpireTime != nil {
		in, out := &in.ExpireTime, &out.ExpireTime
		*out = new(string)
		**out = **in
	}
	if in.Ttl != nil {
		in, out := &in.Ttl, &out.Ttl
		*out = new(string)
		**out = **in
	}
	if in.Etag != nil {
		in, out := &in.Etag, &out.Etag
		*out = new(string)
		**out = **in
	}
	if in.Rotation != nil {
		in, out := &in.Rotation, &out.Rotation
		*out = new(Rotation)
		(*in).DeepCopyInto(*out)
	}
	if in.VersionAliases != nil {
		in, out := &in.VersionAliases, &out.VersionAliases
		*out = make(map[string]int64, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CustomerManagedEncryption != nil {
		in, out := &in.CustomerManagedEncryption, &out.CustomerManagedEncryption
		*out = new(CustomerManagedEncryption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secret.
func (in *Secret) DeepCopy() *Secret {
	if in == nil {
		return nil
	}
	out := new(Secret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretManagerSecret) DeepCopyInto(out *SecretManagerSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretManagerSecret.
func (in *SecretManagerSecret) DeepCopy() *SecretManagerSecret {
	if in == nil {
		return nil
	}
	out := new(SecretManagerSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretManagerSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretManagerSecretList) DeepCopyInto(out *SecretManagerSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecretManagerSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretManagerSecretList.
func (in *SecretManagerSecretList) DeepCopy() *SecretManagerSecretList {
	if in == nil {
		return nil
	}
	out := new(SecretManagerSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretManagerSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretManagerSecretObservedState) DeepCopyInto(out *SecretManagerSecretObservedState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretManagerSecretObservedState.
func (in *SecretManagerSecretObservedState) DeepCopy() *SecretManagerSecretObservedState {
	if in == nil {
		return nil
	}
	out := new(SecretManagerSecretObservedState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretManagerSecretSpec) DeepCopyInto(out *SecretManagerSecretSpec) {
	*out = *in
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
	if in.Replication != nil {
		in, out := &in.Replication, &out.Replication
		*out = new(Replication)
		(*in).DeepCopyInto(*out)
	}
	if in.TopicRefs != nil {
		in, out := &in.TopicRefs, &out.TopicRefs
		*out = make([]TopicRef, len(*in))
		copy(*out, *in)
	}
	if in.ExpireTime != nil {
		in, out := &in.ExpireTime, &out.ExpireTime
		*out = new(string)
		**out = **in
	}
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(string)
		**out = **in
	}
	if in.Rotation != nil {
		in, out := &in.Rotation, &out.Rotation
		*out = new(Rotation)
		(*in).DeepCopyInto(*out)
	}
	if in.VersionAliases != nil {
		in, out := &in.VersionAliases, &out.VersionAliases
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretManagerSecretSpec.
func (in *SecretManagerSecretSpec) DeepCopy() *SecretManagerSecretSpec {
	if in == nil {
		return nil
	}
	out := new(SecretManagerSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretManagerSecretStatus) DeepCopyInto(out *SecretManagerSecretStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1alpha1.Condition, len(*in))
		copy(*out, *in)
	}
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.ExternalRef != nil {
		in, out := &in.ExternalRef, &out.ExternalRef
		*out = new(string)
		**out = **in
	}
	if in.ObservedState != nil {
		in, out := &in.ObservedState, &out.ObservedState
		*out = new(SecretManagerSecretObservedState)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretManagerSecretStatus.
func (in *SecretManagerSecretStatus) DeepCopy() *SecretManagerSecretStatus {
	if in == nil {
		return nil
	}
	out := new(SecretManagerSecretStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Topic) DeepCopyInto(out *Topic) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Topic.
func (in *Topic) DeepCopy() *Topic {
	if in == nil {
		return nil
	}
	out := new(Topic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopicRef) DeepCopyInto(out *TopicRef) {
	*out = *in
	out.PubSubTopicRef = in.PubSubTopicRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopicRef.
func (in *TopicRef) DeepCopy() *TopicRef {
	if in == nil {
		return nil
	}
	out := new(TopicRef)
	in.DeepCopyInto(out)
	return out
}