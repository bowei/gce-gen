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

package main

import (
	"context"
	"flag"
	"log"

	"github.com/golang/glog"
	"golang.org/x/oauth2/google"

	alpha "google.golang.org/api/compute/v0.alpha"
	beta "google.golang.org/api/compute/v0.beta"
	ga "google.golang.org/api/compute/v1"

	"github.com/bowei/gce-gen/pkg/cloud"
	"github.com/bowei/gce-gen/pkg/cloud/filter"
	"github.com/bowei/gce-gen/pkg/cloud/meta"
)

var flags = struct {
	usemock bool
}{}

func init() {
	flag.BoolVar(&flags.usemock, "usemock", false, "usemock")
}

func mockCloud() cloud.Cloud {
	mock := cloud.NewMockGCE()
	mock.MockZones.Objects[*meta.ZonalKey("abc", "us-central1-b")] = &cloud.MockZonesObj{
		ga.Zone{Name: "us-central1-b"},
	}
	return mock
}

func realCloud() cloud.Cloud {
	var g *ga.Service
	var a *alpha.Service
	var b *beta.Service

	c, err := google.DefaultClient(context.Background(), ga.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	g, err = ga.New(c)
	if err != nil {
		log.Fatal(err)
	}
	c, err = google.DefaultClient(context.Background(), alpha.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	a, err = alpha.New(c)
	if err != nil {
		log.Fatal(err)
	}
	c, err = google.DefaultClient(context.Background(), beta.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	b, err = beta.New(c)
	if err != nil {
		log.Fatal(err)
	}

	gce := cloud.NewGCE(&cloud.Service{
		GA:            g,
		Alpha:         a,
		Beta:          b,
		ProjectRouter: &cloud.SingleProjectRouter{ID: "bowei-gke"},
		RateLimiter:   &cloud.NopRateLimiter{},
	})
	return gce
}

func main() {
	flag.Parse()

	var c cloud.Cloud
	if flags.usemock {
		c = mockCloud()
	} else {
		c = realCloud()
	}

	glog.Infof("List addresses")
	addrs, err := c.Addresses().List(context.Background(), "us-central1", filter.None)
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		glog.Infof("addr = %+v\n", addr)
	}

	glog.Infof("List firewalls")
	fws, err := c.Firewalls().List(context.Background(), filter.None)
	if err != nil {
		panic(err)
	}
	for _, fw := range fws {
		glog.Infof("fw = %+v\n", fw)
	}

	// Create, Get and Delete a firewall.
	key := *meta.GlobalKey("abc")
	fw := &ga.Firewall{
		Allowed:      []*ga.FirewallAllowed{{IPProtocol: "tcp", Ports: []string{"80"}}},
		Network:      "projects/bowei-gke/global/networks/custom-empty",
		Direction:    "INGRESS",
		SourceRanges: []string{"104.155.174.199/32"},
	}
	if err := c.Firewalls().Insert(context.Background(), key, fw); err != nil {
		glog.Fatalf("Firewall insert error %v", err)
	}

	glog.Infof("Firewall %v created", key)
	if fw, err = c.Firewalls().Get(context.Background(), key); err != nil {
		glog.Fatalf("Firewall get error %v", err)
	}
	glog.Infof("Firewall is %+v", fw)

	glog.Infof("List firewalls")
	fws, err = c.Firewalls().List(context.Background(), filter.None)
	if err != nil {
		panic(err)
	}
	for _, fw := range fws {
		glog.Infof("fw = %+v\n", fw)
	}

	glog.Infof("List firewalls (starting with a.*)")
	fws, err = c.Firewalls().List(context.Background(), filter.Regexp("name", "a.*"))
	if err != nil {
		panic(err)
	}
	for _, fw := range fws {
		glog.Infof("fw = %+v\n", fw)
	}

	if err := c.Firewalls().Delete(context.Background(), key); err != nil {
		glog.Fatalf("Firewall delete error %v", err)
	}
	glog.Infof("Firewall %v deleted", key)

	if proj, err := c.Projects().Get(context.Background(), "bowei-gke"); err == nil {
		glog.Infof("Project %+v", proj)
		for _, q := range proj.Quotas {
			glog.Infof("Project quota %+v", q)
		}
	} else {
		glog.Errorf("Projects.Get: %v", err)
	}
}
