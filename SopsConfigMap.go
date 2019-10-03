// Copyright 2019 Agilicus Incorporated
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"path/filepath"

	"go.mozilla.org/sops/decrypt"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr                    ifc.Loader
	rf                     *resmap.Factory
	types.GeneratorOptions `json:"generatorOptions,omitempty" yaml:"generatorOptions,omitempty"`
	types.ConfigMapArgs

	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
	// List of keys to use in database lookups
	Keys []string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

// KustomizePlugin ...
//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.ConfigMapArgs = types.ConfigMapArgs{}
	p.GeneratorOptions = types.GeneratorOptions{}
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

// (root string, args []string) (map[string]string, error) {
func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.ConfigMapArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace
	args.GeneratorArgs.Behavior = "merge"

	if len(p.Source) == 0 {
		p.Source = "secrets.enc.yaml"
	}

	secret := make(map[string]string)
	secretFile := filepath.Join(p.ldr.Root(), p.Source)

	v, err := decrypt.File(secretFile, "yaml")
	if err != nil {
		log.Fatalf("error: cannot decode file %s :: %v", secretFile, err)
	}
	err = yaml.Unmarshal([]byte(v), &secret)
	if err != nil {
		log.Fatalf("error: cannot unmarshal %s as yaml :: %v", secretFile, err)
	}

	if len(p.Keys) == 0 {
		for k := range secret {
			p.Keys = append(p.Keys, k)
		}
	}

	for _, k := range p.Keys {
		v, ok := secret[k]
		if ok {
			args.LiteralSources = append(args.LiteralSources, k+"="+v)
		} else {
			log.Fatalf("error: key <%s> not present in %s\n", k, secretFile)
		}
	}

	return p.rf.FromConfigMapArgs(p.ldr, &p.GeneratorOptions, args)
}
