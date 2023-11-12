.PHONY: push

out/undionly.kpxe: build matchbox.ipxe
	./build

push: out/undionly.kpxe
	cp out/undionly.kpxe ../ansible/files/undionly.kpxe
