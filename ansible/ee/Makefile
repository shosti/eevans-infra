.PHONY: build push

build:
	ansible-builder build -v=3 --tag=harbor.eevans.me/library/awx-ee:2.11 --container-runtime=docker

push: build
	docker push harbor.eevans.me/library/awx-ee:2.11
