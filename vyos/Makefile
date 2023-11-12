.PHONY: build push clean

VYOS_VERSION := equuleus
VYOS_COMMIT := 1.3.4

default: build

target/vyos-%.iso:
	bin/build $(VYOS_VERSION) $(VYOS_COMMIT)

target/vyos-%.iso.asc: target/vyos-$(VYOS_VERSION)-$(VYOS_COMMIT).iso
	gpg --armour --detach-sign --default-key F0D9653827AD9C47a --output $@ $<

build: target/vyos-$(VYOS_VERSION)-$(VYOS_COMMIT).iso

push: target/vyos-$(VYOS_VERSION)-$(VYOS_COMMIT).iso target/vyos-$(VYOS_VERSION)-$(VYOS_COMMIT).iso.asc
	bin/push $<

clean:
	sudo rm -rf target build
