// Copyright 2019 Agilicus Incorporated
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"go.mozilla.org/sops/decrypt"
	"log"
	"path/filepath"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr ifc.Loader
	rf  *resmap.Factory
	types.GeneratorOptions
	types.SecretArgs
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
	// List of keys to use in database lookups
	Keys []string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.SecretArgs = types.SecretArgs{}
	p.GeneratorOptions = types.GeneratorOptions{}
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

// (root string, args []string) (map[string]string, error) {
func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace
	//	args.GeneratorArgs.Behavior = "merge"
	log.Printf("args: %+v\n", args)
	log.Printf("\np: %+v\n", p)

	if len(p.Source) == 0 {
		p.Source = "secrets.enc.yaml"
	}

	secret := make(map[string]string)
	secret_file := filepath.Join(p.ldr.Root(), p.Source)

	v, err := decrypt.File(secret_file, "yaml")
	if err != nil {
		log.Fatalf("error: cannot decode file %s :: %v", secret_file, err)
	}
	err = yaml.Unmarshal([]byte(v), &secret)
	if err != nil {
		log.Fatalf("error: cannot unmarshal %s as yaml :: %v", secret_file, err)
	}

	for _, k := range p.Keys {
		v, ok := secret[k]
		if ok {
			args.LiteralSources = append(args.LiteralSources, k+"="+v)
		} else {
			log.Fatalf("error: key <%s> not present in %s\n", k, secret_file)
		}
	}

	//	opts := types.GeneratorOptions{}
	//	opts.Behavior{Behavior: types.GenerationBehavior.BehaviorMerge}
	log.Printf("BEFORE")
	//log.Printf("IDS: %v", opts)

	log.Printf("\n\nargs: %+v\n", args)
	resm, err := p.rf.FromSecretArgs(p.ldr, nil, args)
	log.Printf("\nresm: %+v", resm)
	return p.rf.FromSecretArgs(p.ldr, &p.GeneratorOptions, args)
}
