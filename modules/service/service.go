// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package service

type Service interface {
	Name() string
}

type Generator func() (Service, error)

func Register(key string, gen Generator) {
	if _, ok := services[key]; ok {
		panic("service has been registered")
	}
	services[key] = gen
}

func GetServices() List {
	var (
		l   = make(List)
		err error
	)
	for k, gen := range services {
		if l[k], err = gen(); err != nil {
			panic("service initialize fail")
		}
	}
	return l
}

var services = make(Generators)

type Generators map[string]Generator

type List map[string]Service

func (g List) Get(key string) Service {
	if v, ok := g[key]; ok {
		return v
	}
	panic("service not found")
}

func (g List) GetOrNot(key string) (Service, bool) {
	v, ok := g[key]
	return v, ok
}

func (g List) Add(key string, service Service) {
	if _, ok := g[key]; ok {
		panic("service exist")
	}
	g[key] = service
}
