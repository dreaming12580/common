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

// Package util provides various helper routines.
package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	apiv1 "github.com/kubeflow/common/pkg/apis/common/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/code-generator/pkg/util"
	_ "k8s.io/kube-openapi/pkg/common"
)

const (
	// EnvKubeflowNamespace is a environment variable for namespace when deployed on kubernetes
	EnvKubeflowNamespace = "KUBEFLOW_NAMESPACE"
)

// Pformat returns a pretty format output of any value that can be marshaled to JSON.
func Pformat(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	valueJSON, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Warningf("Couldn't pretty format %v, error: %v", value, err)
		return fmt.Sprintf("%v", value)
	}
	return string(valueJSON)
}

// src is variable initialized with random value.
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString generates a random string of the desired length.
//
// The string is DNS-1035 label compliant; i.e. its only alphanumeric lowercase.
// From: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func VerifyServicesConfig(job metav1.Object, services []*v1.Service, rtype apiv1.ReplicaType, spec *apiv1.ReplicaSpec) error {
	rt := strings.ToLower(string(rtype))
	jobName := job.GetName()

	n := jobName + "-" + strings.ToLower(rt) + "-" + "123"
	serviceName := strings.Replace(n, "/", "-", -1)
	if len(serviceName) > 63 {
		log.Warningf("JobName: %s, length: %d is too long", jobName, len(serviceName))
		return fmt.Errorf("the Job: %s Name is too long", jobName)
	}
	return nil
}

func VerifyJobConfig(job metav1.Object) error {
	jobName := job.GetName()

	// 52 is very conservative value, Because K8s resource obeject metedata.name has Limit of length 63
	// 63 minus length of replica type and job index, remaining length is going to be about 52
	if len(jobName) > 52 {
		log.Warningf("JobName: %s, length: %d is too long", jobName, len(jobName))
		return fmt.Errorf("the Job: %s Name is too long", jobName)
	}
	return nil
}
