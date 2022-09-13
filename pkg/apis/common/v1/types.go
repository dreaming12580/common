// Copyright 2018 The Kubeflow Authors
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
package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
// JobStatus represents the current observed state of the training Job.
type JobStatus struct {
	// Conditions is an array of current observed job conditions.
	Conditions []JobCondition `json:"conditions" protobuf:"bytes,1,rep,name=conditions"`

	// ReplicaStatuses is map of ReplicaType and ReplicaStatus,
	// specifies the status of each replica.
	ReplicaStatuses map[ReplicaType]*ReplicaStatus `json:"replicaStatuses" protobuf:"bytes,2,rep,name=replicaStatuses"`

	// Represents time when the job was acknowledged by the job controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty" protobuf:"bytes,3,opt,name=startTime"`

	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty" protobuf:"bytes,4,opt,name=completionTime"`

	// Represents last time when the job was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty" protobuf:"bytes,5,opt,name=lastReconcileTime"`
}

// +k8s:openapi-gen=true
// ReplicaType represents the type of the replica. Each operator needs to define its
// own set of ReplicaTypes.
type ReplicaType string

// +k8s:openapi-gen=true
// ReplicaStatus represents the current observed state of the replica.
type ReplicaStatus struct {
	// The number of actively running pods.
	Active int32 `json:"active,omitempty" protobuf:"varint,1,opt,name=active"`

	// The number of pods which reached phase Succeeded.
	Succeeded int32 `json:"succeeded,omitempty" protobuf:"varint,2,opt,name=succeed"`

	// The number of pods which reached phase Failed.
	Failed int32 `json:"failed,omitempty" protobuf:"varint,3,opt,name=failed"`

	// A label selector is a label query over a set of resources. The result of matchLabels and
	// matchExpressions are ANDed. An empty label selector matches all objects. A null
	// label selector matches no objects.
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,4,opt,name=labelSelector"`
}

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
// ReplicaSpec is a description of the replica
type ReplicaSpec struct {
	// Replicas is the desired number of replicas of the given template.
	// If unspecified, defaults to 1.
	Replicas *int32 `json:"replicas,omitempty" protobuf:"bytes,1,opt,name=replicas"`

	// Template is the object that describes the pod that
	// will be created for this replica. RestartPolicy in PodTemplateSpec
	// will be overide by RestartPolicy in ReplicaSpec
	Template v1.PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,2,opt,name=template"`

	// Restart policy for all replicas within the job.
	// One of Always, OnFailure, Never and ExitCode.
	// Default to Never.
	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy"`
}

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
// JobCondition describes the state of the job at a certain point.
type JobCondition struct {
	// Type of job condition.
	Type JobConditionType `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,5,opt,name=lastUpdateTime"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,6,opt,name=lastTransitionTime"`
}

// +k8s:openapi-gen=true
// JobConditionType defines all kinds of types of JobStatus.
type JobConditionType string

const (
	// JobCreated means the job has been accepted by the system,
	// but one or more of the pods/services has not been started.
	// This includes time before pods being scheduled and launched.
	JobCreated JobConditionType = "Created"

	// JobRunning means all sub-resources (e.g. services/pods) of this job
	// have been successfully scheduled and launched.
	// The training is running without error.
	JobRunning JobConditionType = "Running"

	// JobRestarting means one or more sub-resources (e.g. services/pods) of this job
	// reached phase failed but maybe restarted according to it's restart policy
	// which specified by user in v1.PodTemplateSpec.
	// The training is freezing/pending.
	JobRestarting JobConditionType = "Restarting"

	// JobSucceeded means all sub-resources (e.g. services/pods) of this job
	// reached phase have terminated in success.
	// The training is complete without error.
	JobSucceeded JobConditionType = "Succeeded"

	// JobFailed means one or more sub-resources (e.g. services/pods) of this job
	// reached phase failed with no restarting.
	// The training has failed its execution.
	JobFailed JobConditionType = "Failed"

	// JobSuspended means sub-resources (e.g. services/pods) of this job
	// has been terminated.
	JobSuspended JobConditionType = "Suspended"

	// JobResumed means job Resumed from suspended
	JobResumed JobConditionType = "Resumed"

	// JobPartialSucceed means all sub-resources (e.g. services/pods) of this job's one worker
	// reached phase have terminated in success.
	JobPartialSucceeded JobConditionType = "PartialSucceeded"
)

// +k8s:openapi-gen=true
// CleanPodPolicy describes how to deal with pods when the job is finished.
type CleanPodPolicy string

const (
	CleanPodPolicyUndefined CleanPodPolicy = ""
	CleanPodPolicyAll       CleanPodPolicy = "All"
	CleanPodPolicyRunning   CleanPodPolicy = "Running"
	CleanPodPolicyNone      CleanPodPolicy = "None"
)

// +k8s:openapi-gen=true
// RestartPolicy describes how the replicas should be restarted.
// Only one of the following restart policies may be specified.
// If none of the following policies is specified, the default one
// is RestartPolicyAlways.
type RestartPolicy string

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"

	// RestartPolicyExitCode policy means that user should add exit code by themselves,
	// The job operator will check these exit codes to
	// determine the behavior when an error occurs:
	// - 1-127: permanent error, do not restart.
	// - 128-255: retryable error, will restart the pod.
	RestartPolicyExitCode RestartPolicy = "ExitCode"
)

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
// RunPolicy encapsulates various runtime policies of the distributed training
// job, for example how to clean up resources and how long the job can stay
// active.
type RunPolicy struct {
	// CleanPodPolicy defines the policy to kill pods after the job completes.
	// Default to Running.
	CleanPodPolicy *CleanPodPolicy `json:"cleanPodPolicy,omitempty" protobuf:"bytes,1,opt,name=cleanPodPolicy"`

	// TTLSecondsAfterFinished is the TTL to clean up jobs.
	// It may take extra ReconcilePeriod seconds for the cleanup, since
	// reconcile gets called periodically.
	// Default to infinite.
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty" protobuf:"bytes,2,opt,name=ttlSecondsAfterFinished"`

	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer.
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" protobuf:"bytes,3,opt,name=activeDeadlineSeconds"`

	// Optional number of retries before marking this job failed.
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty" protobuf:"bytes,4,opt,name=backoffLimit"`

	// SchedulingPolicy defines the policy related to scheduling, e.g. gang-scheduling
	// +optional
	SchedulingPolicy *SchedulingPolicy `json:"schedulingPolicy,omitempty" protobuf:"bytes,5,opt,name=schedulingPolicy"`

	// Suspend specifies whether the Job controller should create Pods or not. If
	// a Job is created with suspend set to true, no Pods are created by the Job
	// controller. If a Job is suspended after creation (i.e. the flag goes from
	// false to true), the Job controller will delete all active Pods associated
	// with this Job. Users must design their workload to gracefully handle this.
	// Suspending a Job will reset the StartTime field of the Job, effectively
	// resetting the ActiveDeadlineSeconds timer too. This is an alpha field and
	// requires the SuspendJob feature gate to be enabled; otherwise this field
	// may not be set to true. Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty" protobuf:"bytes,6,opt,name=suspend"`
}

// +k8s:openapi-gen=true
// SchedulingPolicy encapsulates various scheduling policies of the distributed training
// job, for example `minAvailable` for gang-scheduling.
type SchedulingPolicy struct {
	MinAvailable  *int32           `json:"minAvailable,omitempty" protobuf:"bytes,1,opt,name=minAvailable"`
	Queue         string           `json:"queue,omitempty" protobuf:"bytes,2,opt,name=queue"`
	MinResources  *v1.ResourceList `json:"minResources,omitempty" protobuf:"bytes,3,opt,name=minResources"`
	PriorityClass string           `json:"priorityClass,omitempty" protobuf:"bytes,4,opt,name=priorityClass"`
}
