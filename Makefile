all: ~/bin/kustomize ~/.config/kustomize/plugin/kustomize-sops/v1/sopssecret/SopsSecret.so sigs.k8s.io/kustomize/go.mod

sigs.k8s.io/kustomize/go.mod:
	export GO111MODULE=on
	mkdir -p sigs.k8s.io
	git clone -b fix/sops-integration https://github.com/rosscdh/kustomize.git sigs.k8s.io/kustomize
	# cp -R ../kustomize  sigs.k8s.io/
	
	mkdir -p ~/.config/kustomize/plugin/kustomize-sops/v1/sopssecret
	ln -s $$PWD/SopsSecret.go $$PWD/sigs.k8s.io/kustomize/plugin/
	patch -p1 < kustomize.patch
	patch -p1 < kustomize-enable.patch

~/.config/kustomize/plugin/kustomize-sops/v1/sopssecret/SopsSecret.so: ~/bin/kustomize SopsSecret.go
	(cd sigs.k8s.io/kustomize; go build -buildmode plugin -o ~/.config/kustomize/plugin/kustomize-sops/v1/sopssecret/SopsSecret.so plugin/SopsSecret.go)

~/bin/kustomize: sigs.k8s.io/kustomize/go.mod
	(cd sigs.k8s.io/kustomize; go build  -o ~/bin/kustomize cmd/kustomize/main.go)

