.PHONY: default bootstrap prep-nodes flux-system

default:
	echo "Changes should mostly be applied through flux"

bootstrap:
	bin/bootstrap

prep-nodes:
	bin/prep-nodes

flux-system:
	flux install  --components-extra=image-reflector-controller,image-automation-controller --export > deploy/flux-system/flux-system.yaml
