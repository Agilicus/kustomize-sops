// +build plugin

package main

import (
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
)

type plugin struct{}

var KVSource plugin

func (p plugin) Get(root string, args []string) (map[string]string, error) {

	secret := make(map[string]string)

	secret_file := filepath.Join(root, "secrets.enc.yaml")

	v, err := decrypt.File(secret_file, "yaml")
	if err != nil {
		log.Fatalf("error: cannot decode file %s :: %v", secret_file, err)
	}
	err = yaml.Unmarshal([]byte(v), &secret)
	if err != nil {
		log.Fatalf("error: cannot unmarshal secrets.enc.yaml as yaml in %s :: %v", root, err)
	}

	r := make(map[string]string)

	for _, k := range args {
		v, ok := secret[k]
		if ok {
			r[k] = v
		} else {
			log.Fatalf("error: key <%s> not present in secrets.enc.yaml in %s\n", k, root)
		}
	}

	return r, nil
}
