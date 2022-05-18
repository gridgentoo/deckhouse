//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021 Flant JSC

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

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonParameters) DeepCopyInto(out *CommonParameters) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(string)
		**out = **in
	}
	if in.EvictFailedBarePods != nil {
		in, out := &in.EvictFailedBarePods, &out.EvictFailedBarePods
		*out = new(bool)
		**out = **in
	}
	if in.EvictLocalStoragePods != nil {
		in, out := &in.EvictLocalStoragePods, &out.EvictLocalStoragePods
		*out = new(bool)
		**out = **in
	}
	if in.EvictSystemCriticalPods != nil {
		in, out := &in.EvictSystemCriticalPods, &out.EvictSystemCriticalPods
		*out = new(bool)
		**out = **in
	}
	if in.IgnorePVCPods != nil {
		in, out := &in.IgnorePVCPods, &out.IgnorePVCPods
		*out = new(bool)
		**out = **in
	}
	if in.MaxNoOfPodsToEvictPerNode != nil {
		in, out := &in.MaxNoOfPodsToEvictPerNode, &out.MaxNoOfPodsToEvictPerNode
		*out = new(int)
		**out = **in
	}
	if in.MaxNoOfPodsToEvictPerNamespace != nil {
		in, out := &in.MaxNoOfPodsToEvictPerNamespace, &out.MaxNoOfPodsToEvictPerNamespace
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonParameters.
func (in *CommonParameters) DeepCopy() *CommonParameters {
	if in == nil {
		return nil
	}
	out := new(CommonParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Descheduler) DeepCopyInto(out *Descheduler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Descheduler.
func (in *Descheduler) DeepCopy() *Descheduler {
	if in == nil {
		return nil
	}
	out := new(Descheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Descheduler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeschedulerDeploymentTemplate) DeepCopyInto(out *DeschedulerDeploymentTemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ResourcesRequests != nil {
		in, out := &in.ResourcesRequests, &out.ResourcesRequests
		*out = new(ResourcesRequests)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeschedulerDeploymentTemplate.
func (in *DeschedulerDeploymentTemplate) DeepCopy() *DeschedulerDeploymentTemplate {
	if in == nil {
		return nil
	}
	out := new(DeschedulerDeploymentTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeschedulerPolicy) DeepCopyInto(out *DeschedulerPolicy) {
	*out = *in
	in.CommonParameters.DeepCopyInto(&out.CommonParameters)
	in.Strategies.DeepCopyInto(&out.Strategies)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeschedulerPolicy.
func (in *DeschedulerPolicy) DeepCopy() *DeschedulerPolicy {
	if in == nil {
		return nil
	}
	out := new(DeschedulerPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeschedulerSpec) DeepCopyInto(out *DeschedulerSpec) {
	*out = *in
	in.DeploymentTemplate.DeepCopyInto(&out.DeploymentTemplate)
	in.DeschedulerPolicy.DeepCopyInto(&out.DeschedulerPolicy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeschedulerSpec.
func (in *DeschedulerSpec) DeepCopy() *DeschedulerSpec {
	if in == nil {
		return nil
	}
	out := new(DeschedulerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeschedulerStatus) DeepCopyInto(out *DeschedulerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeschedulerStatus.
func (in *DeschedulerStatus) DeepCopy() *DeschedulerStatus {
	if in == nil {
		return nil
	}
	out := new(DeschedulerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeschedulerStrategies) DeepCopyInto(out *DeschedulerStrategies) {
	*out = *in
	if in.RemoveDuplicates != nil {
		in, out := &in.RemoveDuplicates, &out.RemoveDuplicates
		*out = new(RemoveDuplicates)
		(*in).DeepCopyInto(*out)
	}
	if in.LowNodeUtilization != nil {
		in, out := &in.LowNodeUtilization, &out.LowNodeUtilization
		*out = new(LowNodeUtilization)
		(*in).DeepCopyInto(*out)
	}
	if in.HighNodeUtilization != nil {
		in, out := &in.HighNodeUtilization, &out.HighNodeUtilization
		*out = new(HighNodeUtilization)
		(*in).DeepCopyInto(*out)
	}
	if in.RemovePodsViolatingInterPodAntiAffinity != nil {
		in, out := &in.RemovePodsViolatingInterPodAntiAffinity, &out.RemovePodsViolatingInterPodAntiAffinity
		*out = new(RemovePodsViolatingInterPodAntiAffinity)
		(*in).DeepCopyInto(*out)
	}
	if in.RemovePodsViolatingNodeAffinity != nil {
		in, out := &in.RemovePodsViolatingNodeAffinity, &out.RemovePodsViolatingNodeAffinity
		*out = new(RemovePodsViolatingNodeAffinity)
		(*in).DeepCopyInto(*out)
	}
	if in.RemovePodsViolatingNodeTaints != nil {
		in, out := &in.RemovePodsViolatingNodeTaints, &out.RemovePodsViolatingNodeTaints
		*out = new(RemovePodsViolatingNodeTaints)
		(*in).DeepCopyInto(*out)
	}
	if in.RemovePodsViolatingTopologySpreadConstraint != nil {
		in, out := &in.RemovePodsViolatingTopologySpreadConstraint, &out.RemovePodsViolatingTopologySpreadConstraint
		*out = new(RemovePodsViolatingTopologySpreadConstraint)
		(*in).DeepCopyInto(*out)
	}
	if in.RemovePodsHavingTooManyRestarts != nil {
		in, out := &in.RemovePodsHavingTooManyRestarts, &out.RemovePodsHavingTooManyRestarts
		*out = new(RemovePodsHavingTooManyRestarts)
		(*in).DeepCopyInto(*out)
	}
	if in.PodLifeTime != nil {
		in, out := &in.PodLifeTime, &out.PodLifeTime
		*out = new(PodLifeTime)
		(*in).DeepCopyInto(*out)
	}
	if in.RemoveFailedPods != nil {
		in, out := &in.RemoveFailedPods, &out.RemoveFailedPods
		*out = new(RemoveFailedPods)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeschedulerStrategies.
func (in *DeschedulerStrategies) DeepCopy() *DeschedulerStrategies {
	if in == nil {
		return nil
	}
	out := new(DeschedulerStrategies)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HighNodeUtilization) DeepCopyInto(out *HighNodeUtilization) {
	*out = *in
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.NodeResourceUtilizationThresholds != nil {
		in, out := &in.NodeResourceUtilizationThresholds, &out.NodeResourceUtilizationThresholds
		*out = new(NodeResourceUtilizationThresholdsFiltering)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HighNodeUtilization.
func (in *HighNodeUtilization) DeepCopy() *HighNodeUtilization {
	if in == nil {
		return nil
	}
	out := new(HighNodeUtilization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelSelectorFiltering) DeepCopyInto(out *LabelSelectorFiltering) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelSelectorFiltering.
func (in *LabelSelectorFiltering) DeepCopy() *LabelSelectorFiltering {
	if in == nil {
		return nil
	}
	out := new(LabelSelectorFiltering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LowNodeUtilization) DeepCopyInto(out *LowNodeUtilization) {
	*out = *in
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.NodeResourceUtilizationThresholds != nil {
		in, out := &in.NodeResourceUtilizationThresholds, &out.NodeResourceUtilizationThresholds
		*out = new(NodeResourceUtilizationThresholdsFiltering)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LowNodeUtilization.
func (in *LowNodeUtilization) DeepCopy() *LowNodeUtilization {
	if in == nil {
		return nil
	}
	out := new(LowNodeUtilization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Namespaces) DeepCopyInto(out *Namespaces) {
	*out = *in
	if in.Include != nil {
		in, out := &in.Include, &out.Include
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Namespaces.
func (in *Namespaces) DeepCopy() *Namespaces {
	if in == nil {
		return nil
	}
	out := new(Namespaces)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacesFiltering) DeepCopyInto(out *NamespacesFiltering) {
	*out = *in
	in.Namespaces.DeepCopyInto(&out.Namespaces)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacesFiltering.
func (in *NamespacesFiltering) DeepCopy() *NamespacesFiltering {
	if in == nil {
		return nil
	}
	out := new(NamespacesFiltering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeFitFiltering) DeepCopyInto(out *NodeFitFiltering) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeFitFiltering.
func (in *NodeFitFiltering) DeepCopy() *NodeFitFiltering {
	if in == nil {
		return nil
	}
	out := new(NodeFitFiltering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeResourceUtilizationThresholdsFiltering) DeepCopyInto(out *NodeResourceUtilizationThresholdsFiltering) {
	*out = *in
	if in.Thresholds != nil {
		in, out := &in.Thresholds, &out.Thresholds
		*out = make(ResourceThresholds, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TargetThresholds != nil {
		in, out := &in.TargetThresholds, &out.TargetThresholds
		*out = make(ResourceThresholds, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeResourceUtilizationThresholdsFiltering.
func (in *NodeResourceUtilizationThresholdsFiltering) DeepCopy() *NodeResourceUtilizationThresholdsFiltering {
	if in == nil {
		return nil
	}
	out := new(NodeResourceUtilizationThresholdsFiltering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLifeTime) DeepCopyInto(out *PodLifeTime) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	if in.PodLifeTime != nil {
		in, out := &in.PodLifeTime, &out.PodLifeTime
		*out = new(PodLifeTimeParameters)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLifeTime.
func (in *PodLifeTime) DeepCopy() *PodLifeTime {
	if in == nil {
		return nil
	}
	out := new(PodLifeTime)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLifeTimeParameters) DeepCopyInto(out *PodLifeTimeParameters) {
	*out = *in
	if in.MaxPodLifeTimeSeconds != nil {
		in, out := &in.MaxPodLifeTimeSeconds, &out.MaxPodLifeTimeSeconds
		*out = new(uint)
		**out = **in
	}
	if in.PodStatusPhases != nil {
		in, out := &in.PodStatusPhases, &out.PodStatusPhases
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLifeTimeParameters.
func (in *PodLifeTimeParameters) DeepCopy() *PodLifeTimeParameters {
	if in == nil {
		return nil
	}
	out := new(PodLifeTimeParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodsHavingTooManyRestartsParameters) DeepCopyInto(out *PodsHavingTooManyRestartsParameters) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodsHavingTooManyRestartsParameters.
func (in *PodsHavingTooManyRestartsParameters) DeepCopy() *PodsHavingTooManyRestartsParameters {
	if in == nil {
		return nil
	}
	out := new(PodsHavingTooManyRestartsParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveDuplicates) DeepCopyInto(out *RemoveDuplicates) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.RemoveDuplicates != nil {
		in, out := &in.RemoveDuplicates, &out.RemoveDuplicates
		*out = new(RemoveDuplicatesParameters)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveDuplicates.
func (in *RemoveDuplicates) DeepCopy() *RemoveDuplicates {
	if in == nil {
		return nil
	}
	out := new(RemoveDuplicates)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveDuplicatesParameters) DeepCopyInto(out *RemoveDuplicatesParameters) {
	*out = *in
	if in.ExcludeOwnerKinds != nil {
		in, out := &in.ExcludeOwnerKinds, &out.ExcludeOwnerKinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveDuplicatesParameters.
func (in *RemoveDuplicatesParameters) DeepCopy() *RemoveDuplicatesParameters {
	if in == nil {
		return nil
	}
	out := new(RemoveDuplicatesParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveFailedPods) DeepCopyInto(out *RemoveFailedPods) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.RemoveFailedPods != nil {
		in, out := &in.RemoveFailedPods, &out.RemoveFailedPods
		*out = new(RemoveFailedPodsParameters)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveFailedPods.
func (in *RemoveFailedPods) DeepCopy() *RemoveFailedPods {
	if in == nil {
		return nil
	}
	out := new(RemoveFailedPods)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoveFailedPodsParameters) DeepCopyInto(out *RemoveFailedPodsParameters) {
	*out = *in
	if in.ExcludeOwnerKinds != nil {
		in, out := &in.ExcludeOwnerKinds, &out.ExcludeOwnerKinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MinPodLifetimeSeconds != nil {
		in, out := &in.MinPodLifetimeSeconds, &out.MinPodLifetimeSeconds
		*out = new(uint)
		**out = **in
	}
	if in.Reasons != nil {
		in, out := &in.Reasons, &out.Reasons
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoveFailedPodsParameters.
func (in *RemoveFailedPodsParameters) DeepCopy() *RemoveFailedPodsParameters {
	if in == nil {
		return nil
	}
	out := new(RemoveFailedPodsParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemovePodsHavingTooManyRestarts) DeepCopyInto(out *RemovePodsHavingTooManyRestarts) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.PodsHavingTooManyRestarts != nil {
		in, out := &in.PodsHavingTooManyRestarts, &out.PodsHavingTooManyRestarts
		*out = new(PodsHavingTooManyRestartsParameters)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemovePodsHavingTooManyRestarts.
func (in *RemovePodsHavingTooManyRestarts) DeepCopy() *RemovePodsHavingTooManyRestarts {
	if in == nil {
		return nil
	}
	out := new(RemovePodsHavingTooManyRestarts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemovePodsViolatingInterPodAntiAffinity) DeepCopyInto(out *RemovePodsViolatingInterPodAntiAffinity) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemovePodsViolatingInterPodAntiAffinity.
func (in *RemovePodsViolatingInterPodAntiAffinity) DeepCopy() *RemovePodsViolatingInterPodAntiAffinity {
	if in == nil {
		return nil
	}
	out := new(RemovePodsViolatingInterPodAntiAffinity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemovePodsViolatingNodeAffinity) DeepCopyInto(out *RemovePodsViolatingNodeAffinity) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.NodeAffinityType != nil {
		in, out := &in.NodeAffinityType, &out.NodeAffinityType
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemovePodsViolatingNodeAffinity.
func (in *RemovePodsViolatingNodeAffinity) DeepCopy() *RemovePodsViolatingNodeAffinity {
	if in == nil {
		return nil
	}
	out := new(RemovePodsViolatingNodeAffinity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemovePodsViolatingNodeTaints) DeepCopyInto(out *RemovePodsViolatingNodeTaints) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
	if in.ExcludedTaints != nil {
		in, out := &in.ExcludedTaints, &out.ExcludedTaints
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemovePodsViolatingNodeTaints.
func (in *RemovePodsViolatingNodeTaints) DeepCopy() *RemovePodsViolatingNodeTaints {
	if in == nil {
		return nil
	}
	out := new(RemovePodsViolatingNodeTaints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemovePodsViolatingTopologySpreadConstraint) DeepCopyInto(out *RemovePodsViolatingTopologySpreadConstraint) {
	*out = *in
	in.ThresholdPrioritiesFiltering.DeepCopyInto(&out.ThresholdPrioritiesFiltering)
	in.NamespacesFiltering.DeepCopyInto(&out.NamespacesFiltering)
	in.LabelSelectorFiltering.DeepCopyInto(&out.LabelSelectorFiltering)
	out.NodeFitFiltering = in.NodeFitFiltering
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemovePodsViolatingTopologySpreadConstraint.
func (in *RemovePodsViolatingTopologySpreadConstraint) DeepCopy() *RemovePodsViolatingTopologySpreadConstraint {
	if in == nil {
		return nil
	}
	out := new(RemovePodsViolatingTopologySpreadConstraint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in ResourceThresholds) DeepCopyInto(out *ResourceThresholds) {
	{
		in := &in
		*out = make(ResourceThresholds, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceThresholds.
func (in ResourceThresholds) DeepCopy() ResourceThresholds {
	if in == nil {
		return nil
	}
	out := new(ResourceThresholds)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcesRequests) DeepCopyInto(out *ResourcesRequests) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcesRequests.
func (in *ResourcesRequests) DeepCopy() *ResourcesRequests {
	if in == nil {
		return nil
	}
	out := new(ResourcesRequests)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ThresholdPrioritiesFiltering) DeepCopyInto(out *ThresholdPrioritiesFiltering) {
	*out = *in
	if in.ThresholdPriority != nil {
		in, out := &in.ThresholdPriority, &out.ThresholdPriority
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ThresholdPrioritiesFiltering.
func (in *ThresholdPrioritiesFiltering) DeepCopy() *ThresholdPrioritiesFiltering {
	if in == nil {
		return nil
	}
	out := new(ThresholdPrioritiesFiltering)
	in.DeepCopyInto(out)
	return out
}
