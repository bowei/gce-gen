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
	"reflect"

	alpha "google.golang.org/api/compute/v0.alpha"
	beta "google.golang.org/api/compute/v0.beta"
	ga "google.golang.org/api/compute/v1"
)

type Version string

const (
	// ReadOnly specifies that the given resource is read-only and should not
	// have insert() or delete() methods generated for the wrapper.
	ReadOnly = 1 << iota
	// CustomOps specifies that an empty interface xxxOps will be generated to
	// enable custom method calls to be attached to the generated service
	// interface.
	CustomOps = 1 << iota

	// VersionGA is the API version in compute.v1.
	VersionGA Version = "ga"
	// VersionAlpha is the API version in computer.v0.alpha.
	VersionAlpha Version = "alpha"
	// VersionBeta is the API version in computer.v0.beta.
	VersionBeta Version = "beta"
)

// AllVersions is a list of all versions of the GCE API.
var AllVersions = []Version{
	VersionGA,
	VersionAlpha,
	VersionBeta,
}

// AllServices are a list of all the services to generate code for. Keep
// this list in lexiographical order by object type.
var AllServices = []*ServiceInfo{
	&ServiceInfo{
		Object:      "Address",
		Service:     "Addresses",
		keyType:     Regional,
		serviceType: reflect.TypeOf(&ga.AddressesService{}),
	},
	&ServiceInfo{
		Object:      "Address",
		Service:     "Addresses",
		version:     VersionAlpha,
		keyType:     Regional,
		serviceType: reflect.TypeOf(&alpha.AddressesService{}),
	},
	&ServiceInfo{
		Object:      "Address",
		Service:     "Addresses",
		version:     VersionBeta,
		keyType:     Regional,
		serviceType: reflect.TypeOf(&beta.AddressesService{}),
	},
	&ServiceInfo{
		Object:      "Address",
		Service:     "GlobalAddresses",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.GlobalAddressesService{}),
	},
	&ServiceInfo{
		Object:      "BackendService",
		Service:     "BackendServices",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.BackendServicesService{}),
		additionalMethods: []string{
			"GetHealth",
			"Update",
		},
	},
	&ServiceInfo{
		Object:            "BackendService",
		Service:           "BackendServices",
		version:           VersionAlpha,
		keyType:           Global,
		serviceType:       reflect.TypeOf(&alpha.BackendServicesService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:      "BackendService",
		Service:     "RegionBackendServices",
		version:     VersionAlpha,
		keyType:     Regional,
		serviceType: reflect.TypeOf(&alpha.RegionBackendServicesService{}),
		additionalMethods: []string{
			"GetHealth",
			"Update",
		},
	},
	&ServiceInfo{
		Object:      "Disk",
		Service:     "Disks",
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&ga.DisksService{}),
	},
	&ServiceInfo{
		Object:      "Disk",
		Service:     "Disks",
		version:     VersionAlpha,
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&alpha.DisksService{}),
	},
	&ServiceInfo{
		Object:      "Disk",
		Service:     "RegionDisks",
		version:     VersionAlpha,
		keyType:     Regional,
		serviceType: reflect.TypeOf(&alpha.DisksService{}),
	},
	&ServiceInfo{
		Object:            "Firewall",
		Service:           "Firewalls",
		keyType:           Global,
		serviceType:       reflect.TypeOf(&ga.FirewallsService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:      "ForwardingRule",
		Service:     "ForwardingRules",
		keyType:     Regional,
		serviceType: reflect.TypeOf(&ga.ForwardingRulesService{}),
	},
	&ServiceInfo{
		Object:      "ForwardingRule",
		Service:     "ForwardingRules",
		version:     VersionAlpha,
		keyType:     Regional,
		serviceType: reflect.TypeOf(&alpha.ForwardingRulesService{}),
	},
	&ServiceInfo{
		Object:            "ForwardingRule",
		Service:           "GlobalForwardingRules",
		keyType:           Global,
		serviceType:       reflect.TypeOf(&ga.GlobalForwardingRulesService{}),
		additionalMethods: []string{"SetTarget"},
	},
	&ServiceInfo{
		Object:            "HealthCheck",
		Service:           "HealthChecks",
		keyType:           Global,
		serviceType:       reflect.TypeOf(&ga.HealthChecksService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:            "HealthCheck",
		Service:           "HealthChecks",
		version:           VersionAlpha,
		keyType:           Global,
		serviceType:       reflect.TypeOf(&alpha.HealthChecksService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:            "HttpHealthCheck",
		Service:           "HttpHealthChecks",
		keyType:           Global,
		serviceType:       reflect.TypeOf(&ga.HttpHealthChecksService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:            "HttpsHealthCheck",
		Service:           "HttpsHealthChecks",
		keyType:           Global,
		serviceType:       reflect.TypeOf(&ga.HttpsHealthChecksService{}),
		additionalMethods: []string{"Update"},
	},
	&ServiceInfo{
		Object:      "InstanceGroup",
		Service:     "InstanceGroups",
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&ga.InstanceGroupsService{}),
		additionalMethods: []string{
			"AddInstances",
			"ListInstances",
			"RemoveInstances",
			"SetNamedPorts",
		},
	},
	&ServiceInfo{
		Object:      "Instance",
		Service:     "Instances",
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&ga.InstancesService{}),
		additionalMethods: []string{
			"AttachDisk",
			"DetachDisk",
		},
	},
	&ServiceInfo{
		Object:      "Instance",
		Service:     "Instances",
		version:     VersionBeta,
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&beta.InstancesService{}),
		additionalMethods: []string{
			"AttachDisk",
			"DetachDisk",
		},
	},
	&ServiceInfo{
		Object:      "Instance",
		Service:     "Instances",
		version:     VersionAlpha,
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&alpha.InstancesService{}),
		additionalMethods: []string{
			"AttachDisk",
			"DetachDisk",
			"UpdateNetworkInterface",
		},
	},
	&ServiceInfo{
		Object:      "NetworkEndpointGroup",
		Service:     "NetworkEndpointGroups",
		version:     VersionAlpha,
		keyType:     Zonal,
		serviceType: reflect.TypeOf(&alpha.NetworkEndpointGroupsService{}),
		additionalMethods: []string{
			"AttachNetworkEndpoints",
			"DetachNetworkEndpoints",
		},
	},
	&ServiceInfo{
		Object:      "Region",
		Service:     "Regions",
		keyType:     Global,
		options:     ReadOnly,
		serviceType: reflect.TypeOf(&ga.RegionsService{}),
	},
	&ServiceInfo{
		Object:      "Route",
		Service:     "Routes",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.RoutesService{}),
	},
	&ServiceInfo{
		Object:      "SslCertificate",
		Service:     "SslCertificates",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.SslCertificatesService{}),
	},
	&ServiceInfo{
		Object:      "TargetHttpProxy",
		Service:     "TargetHttpProxies",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.TargetHttpProxiesService{}),
		additionalMethods: []string{
			"SetUrlMap",
		},
	},
	&ServiceInfo{
		Object:      "TargetHttpsProxy",
		Service:     "TargetHttpsProxies",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.TargetHttpsProxiesService{}),
		additionalMethods: []string{
			"SetSslCertificates",
			"SetUrlMap",
		},
	},
	&ServiceInfo{
		Object:      "TargetPool",
		Service:     "TargetPools",
		keyType:     Regional,
		serviceType: reflect.TypeOf(&ga.TargetPoolsService{}),
		additionalMethods: []string{
			"AddInstance",
			"RemoveInstance",
		},
	},
	&ServiceInfo{
		Object:      "UrlMap",
		Service:     "UrlMaps",
		keyType:     Global,
		serviceType: reflect.TypeOf(&ga.UrlMapsService{}),
		additionalMethods: []string{
			"Update",
		},
	},
	&ServiceInfo{
		Object:      "Zone",
		Service:     "Zones",
		keyType:     Global,
		options:     ReadOnly,
		serviceType: reflect.TypeOf(&ga.ZonesService{}),
	},
}
