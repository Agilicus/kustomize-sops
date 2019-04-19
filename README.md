## kustomize-sops

This is a *very* simple wrapper to allow use of [sops](https://github.com/mozilla/sops)
encoded secrets within [kustomize](https://github.com/kubernetes-sigs/kustomize).

It assumes that there exists a single `secrets.enc.yaml` file, and in it there is a
1-deep YAML representation of `SECRET: VALUE`.

Assume you had a _raw_ secrets as secrets.yaml:
```
CAT: ferocious
DOG: tame
```

You would then encrypt it something like:
```
sops --encrypt --gcp-kms projects/MYPROJECT/locations/global/keyRings/sops/cryptoKeys/sops-key secrets.yaml > secrets.enc.yaml
```

You would use a `kustomization.yaml` file as:

```
secretGenerator:
- name: mysecrets
  kvSources:
  - name: kustomize-sops
    pluginType: go
    args:
    - CAT
    - DOG
```

And then running `kustomize --enable_alpha_goplugins_accept_panic_risk build .` would yield:

```
apiVersion: v1
data:
  CAT: ZmVyb2Npb3Vz
  DOG: dGFtZQ==
kind: Secret
metadata:
  name: mysecrets-4g7fk45c8c
type: Opaque
```

### Install Pre-requisites

```
go get -u github.com/kubernetes-sigs/kustomize
go get -u go.mozilla.org/sops/cmd/sops
```

### Build & Install plugin

```
mkdir -p ~/.config/kustomize/plugins/kvSources
go build -buildmode plugin -o ~/.config/kustomize/plugins/kvSources/kustomize-sops.so kustomize-sops.go
```

### Test/Run

```
kustomize --enable_alpha_goplugins_accept_panic_risk build .
```

### Setup encrypted secrets

```
gcloud auth application-default login
gcloud kms keyrings create sops --location global
gcloud kms keys create sops-key --location global --keyring sops --purpose encryption
gcloud kms keys list --location global --keyring sops
# NAME                                                                      PURPOSE          LABELS  PRIMARY_ID  PRIMARY_STATE
# projects/MYPROJECT/locations/global/keyRings/sops/cryptoKeys/sops-key  ENCRYPT_DECRYPT          1           ENABLED

sops --encrypt --gcp-kms projects/MYPROJECT/locations/global/keyRings/sops/cryptoKeys/sops-key secrets.yaml > secrets.enc.yaml
```
