/*
Copyright 2017 The Kubernetes Authors.

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

package meta

import (
	"fmt"
	"reflect"
)

// ServiceInfo defines the entry for a Service that code will be generated for.
type ServiceInfo struct {
	Object  string
	Service string
	// version if unspecified will be assumed to be VersionGA.
	version     Version
	keyType     KeyType
	serviceType reflect.Type

	additionalMethods []string
	options           int
}

// Version returns the version of the Service, defaulting to GA if APIVersion
// is empty.
func (i *ServiceInfo) Version() Version {
	if i.version == "" {
		return VersionGA
	}
	return i.version
}

func (i *ServiceInfo) VersionField() string {
	switch i.Version() {
	case VersionGA:
		return "GA"
	case VersionAlpha:
		return "Alpha"
	case VersionBeta:
		return "Beta"
	}
	panic(fmt.Errorf("invalid version %q", i.Version()))
}

// WrapType is the name of the wrapper service type.
func (i *ServiceInfo) WrapType() string {
	switch i.Version() {
	case VersionGA:
		return i.Service
	case VersionAlpha:
		return "Alpha" + i.Service
	case VersionBeta:
		return "Beta" + i.Service
	}
	return "Invalid"
}

// WrapTypeOps is the name of the additional operations type.
func (i *ServiceInfo) WrapTypeOps() string {
	return i.WrapType() + "Ops"
}

// FQObjectType is fully qualified name of the object (e.g. compute.Instance).
func (i *ServiceInfo) FQObjectType() string {
	return fmt.Sprintf("%v.%v", i.Version(), i.Object)
}

// ObjectListType is the compute List type for the object (contains Items field).
func (i *ServiceInfo) ObjectListType() string {
	return fmt.Sprintf("%v.%vList", i.Version(), i.Object)
}

// MockWrapType is the name of the concrete mock for this type.
func (i *ServiceInfo) MockWrapType() string {
	return "Mock" + i.WrapType()
}

// MockField is the name of the field in the mock struct.
func (i *ServiceInfo) MockField() string {
	return "Mock" + i.WrapType()
}

// GCEWrapType is the name of the GCE wrapper type.
func (i *ServiceInfo) GCEWrapType() string {
	return "GCE" + i.WrapType()
}

// Field is the name of the GCE struct.
func (i *ServiceInfo) Field() string {
	return "gce" + i.WrapType()
}

func (i *ServiceInfo) Methods() []*method {
	methods := map[string]bool{}
	for _, m := range i.additionalMethods {
		methods[m] = true
	}

	var ret []*method
	for j := 0; j < i.serviceType.NumMethod(); j++ {
		m := i.serviceType.Method(j)
		if _, ok := methods[m.Name]; !ok {
			continue
		}
		sm := newMethod(i, m)
		sm.sanityCheck()
		ret = append(ret, sm)
		methods[m.Name] = false
	}

	for k, b := range methods {
		if b {
			panic(fmt.Errorf("method %q was not found in service %q", k, i.Service))
		}
	}

	return ret
}

// KeyIsGlobal is true if the key is global.
func (i *ServiceInfo) KeyIsGlobal() bool {
	return i.keyType == Global
}

// KeyIsRegional is true if the key is regional.
func (i *ServiceInfo) KeyIsRegional() bool {
	return i.keyType == Regional
}

// KeyIsZonal is true if the key is zonal.
func (i *ServiceInfo) KeyIsZonal() bool {
	return i.keyType == Zonal
}

// GenerateMutations is true if we should generate mutations for the object,
// i.e. insert() and delete().
func (i *ServiceInfo) GenerateMutations() bool {
	return i.options&ReadOnly == 0
}

// GenerateCustomOps is true if we should generated a xxxOps interface for
// adding additional methods to the generated interface.
func (i *ServiceInfo) GenerateCustomOps() bool {
	return i.options&CustomOps != 0
}
