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
	rf        *resmap.Factory
	ldr       ifc.Loader
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// List of keys to use in database lookups
	Keys []string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

// (root string, args []string) (map[string]string, error) {
func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace

	secret := make(map[string]string)
	secret_file := filepath.Join(p.ldr.Root(), "secrets.enc.yaml")

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

	return p.rf.FromSecretArgs(p.ldr, nil, args)
}
