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

// Package callbacks implements a simple list of callbacks to be used when
// creating K8s watchers. The intent is to allow the K8s watchers to
// consolidate all the event handler callbacks from various subsystems into one
// place.
package callbacks

import (
	"github.com/cilium/cilium/pkg/lock"

	"k8s.io/client-go/tools/cache"
)

// New creates a new callbacks list.
func New() *List {
	return &List{
		cbs: make([]cache.ResourceEventHandlerFuncs, 0),
	}
}

// Register registers an event handler's callbacks into the list.
func (q *List) Register(cb cache.ResourceEventHandler) {
	q.Lock()
	q.cbs = append(q.cbs, cache.ResourceEventHandlerFuncs{
		AddFunc:    cb.OnAdd,
		UpdateFunc: cb.OnUpdate,
		DeleteFunc: cb.OnDelete,
	})
	q.Unlock()
}

// Callbacks returns all the event handler callbacks from the list.
func (q *List) Callbacks() []cache.ResourceEventHandlerFuncs {
	q.RLock()
	defer q.RUnlock()
	return q.cbs

}

// List holds event handler callbacks used when reacting to K8s resource /
// object changes in the K8s watchers.
type List struct {
	lock.RWMutex

	cbs []cache.ResourceEventHandlerFuncs
}
