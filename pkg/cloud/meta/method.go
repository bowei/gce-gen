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
	"strings"
)

func newArg(t reflect.Type) *arg {
	ret := &arg{}
Loop:
	for {
		switch t.Kind() {
		case reflect.Ptr:
			ret.numPtr++
			t = t.Elem()
		default:
			ret.pkg = t.PkgPath()
			ret.typeName += t.Name()
			break Loop
		}
	}
	return ret
}

type arg struct {
	pkg, typeName string
	numPtr        int
}

func (a *arg) normalizedPkg() string {
	if a.pkg == "" {
		return ""
	}

	// XXX/bowei -- this is hugely ugly
	parts := strings.Split(a.pkg, "/")
	// Remove vendor prefix.
	for i := 0; i < len(parts); i++ {
		if parts[i] == "vendor" {
			parts = parts[i+1:]
			break
		}
	}
	switch strings.Join(parts, "/") {
	case "google.golang.org/api/compute/v1":
		return "ga."
	case "google.golang.org/api/compute/v0.alpha":
		return "alpha."
	case "google.golang.org/api/compute/v0.beta":
		return "beta."
	default:
		panic(fmt.Errorf("unhandled package %q", a.pkg))
	}
}

func (a *arg) String() string {
	var ret string
	for i := 0; i < a.numPtr; i++ {
		ret += "*"
	}
	ret += a.normalizedPkg()
	ret += a.typeName
	return ret
}

// method is used to generate the calling code non-standard methods.
type method struct {
	*ServiceInfo
	m reflect.Method
}

// argsSkip is the number of arguments to skip when generating the
// synthesized method.
func (mr *method) argsSkip() int {
	switch mr.keyType {
	case Zonal:
		return 4
	case Regional:
		return 4
	case Global:
		return 3
	}
	panic(fmt.Errorf("invalid KeyType %v", mr.keyType))
}

// args return a list of arguments to the method, skipping the first skip
// elements. If nameArgs is true, then the arguments will include a generated
// parameter name (arg<N>). prefix will be added to the parameters.
func (mr *method) args(skip int, nameArgs bool, prefix []string) []string {
	var args []*arg
	fType := mr.m.Func.Type()
	for i := 0; i < fType.NumIn(); i++ {
		t := fType.In(i)
		args = append(args, newArg(t))
	}

	var a []string
	for i := skip; i < fType.NumIn(); i++ {
		if nameArgs {
			a = append(a, fmt.Sprintf("arg%d %s", i-skip, args[i]))
		} else {
			a = append(a, args[i].String())
		}
	}
	return append(prefix, a...)
}

func (mr *method) sanityCheck() {
	fType := mr.m.Func.Type()
	if fType.NumIn() < mr.argsSkip() {
		panic(fmt.Errorf("method %q in service %q, arity = %d which is less than required for auto-generation", mr.Name(), mr.Service, fType.NumIn()))
	}
	for i := 0; i < fType.NumIn(); i++ {
	}
}

func (mr *method) Name() string {
	return mr.m.Name
}

func (mr *method) CallArgs() string {
	var args []string
	for i := mr.argsSkip(); i < mr.m.Func.Type().NumIn(); i++ {
		args = append(args, fmt.Sprintf("arg%d", i-mr.argsSkip()))
	}
	if len(args) == 0 {
		return ""
	}
	return fmt.Sprintf(", %s", strings.Join(args, ", "))
}

func (mr *method) MockHookName() string {
	return mr.m.Name + "Hook"
}

func (mr *method) MockHook() string {
	args := mr.args(mr.argsSkip(), false, []string{
		fmt.Sprintf("*%s", mr.MockWrapType()),
		"context.Context",
		"meta.Key",
	})
	return fmt.Sprintf("%v func(%v) error", mr.MockHookName(), strings.Join(args, ", "))
}

func (mr *method) MockFcnArgs() string {
	args := mr.args(mr.argsSkip(), true, []string{
		"ctx context.Context",
		"key meta.Key",
	})
	return fmt.Sprintf("%v(%v) error", mr.m.Name, strings.Join(args, ", "))
}

func (mr *method) InterfaceFunc() string {
	args := mr.args(mr.argsSkip(), false, []string{"context.Context", "meta.Key"})
	return fmt.Sprintf("%v(%v) error", mr.m.Name, strings.Join(args, ", "))
}
