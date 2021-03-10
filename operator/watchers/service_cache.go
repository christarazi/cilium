// Copyright 2021 Authors of Cilium
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

package watchers

import (
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/logging/logfields"
)

func newServiceCacheCallback(swgSvcs, swgEps *lock.StoppableWaitGroup) *serviceCacheCallback {
	return &serviceCacheCallback{
		swgSvcs: swgSvcs,
		swgEps:  swgEps,
	}
}

// serviceCacheCallback represents an object containing K8s event reaction
// logic for the K8sSvcCache. It implements cache.ResourceEventHandler and is
// used in the watcher.
type serviceCacheCallback struct {
	swgSvcs, swgEps *lock.StoppableWaitGroup
}

func (c *serviceCacheCallback) OnAdd(obj interface{}) {
	svc, ok := obj.(*slim_corev1.Service)
	if !ok {
		return
	}
	log.WithField(logfields.ServiceName, svc.Name).Debugf("Received service addition %+v", svc)
	K8sSvcCache.UpdateService(svc, c.swgSvcs)
}
func (c *serviceCacheCallback) OnUpdate(oldObj, newObj interface{}) {
	svc, ok := newObj.(*slim_corev1.Service)
	if !ok {
		return
	}
	log.WithField(logfields.ServiceName, svc.Name).Debugf("Received service update %+v", svc)
	K8sSvcCache.UpdateService(svc, c.swgSvcs)
}
func (c *serviceCacheCallback) OnDelete(obj interface{}) {
	svc, ok := obj.(*slim_corev1.Service)
	if !ok {
		return
	}
	log.WithField(logfields.ServiceName, svc.Name).Debugf("Received service deletion %+v", svc)
	K8sSvcCache.DeleteService(svc, c.swgSvcs)
}
