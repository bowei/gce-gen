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

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bowei/gce-gen/pkg/cloud/meta"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// ProjectOps is the manually implemented methods for the Projects service.
type ProjectsOps interface {
	Get(ctx context.Context, projectID string) (*compute.Project, error)
}

func (m *MockProjects) Get(ctx context.Context, projectID string) (*compute.Project, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if p, ok := m.Objects[*meta.GlobalKey(projectID)]; ok {
		return p, nil
	}
	return nil, &googleapi.Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("MockProjects %v not found", projectID),
	}
}

func (g *GCEProjects) Get(ctx context.Context, projectID string) (*compute.Project, error) {
	rk := &RateLimitKey{
		ProjectID: projectID,
		Operation: "Get",
		Version:   meta.Version("ga"),
		Service:   "Projects",
	}
	g.s.RateLimiter.Accept(ctx, rk)
	call := g.s.GA.Projects.Get(projectID)
	call.Context(ctx)
	return call.Do()
}
