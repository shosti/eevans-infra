# My Personal Infra Repo

This is my infrastructure-as-code repo for my homelab and other self-hosted
stuff. I'm using [gitslice](https://github.com/GitStartHQ/gitslice/) to
partially mirror an open source version, since the full version has some
encrypted secrets and networking details I'm not comfortable sharing.

Most things won't work out-of-the-box, but hopefully there's something useful if
you're interested in this sort of stuff!

## Quick tour

At a high level, the setup involves:

- A router running [vyos](https://vyos.io/), managed by Ansible.
- A couple of Raspberry Pis (the `netsvc` boxes) that run DNS and some other
  networking-related things, also managed through Ansible.
- A bare-metal Kubernetes cluster that's a bit more exotic:
  - The nodes run [Flatcar Linux](https://www.flatcar-linux.org/), which are
    netbooted through iPXE with [Matchbox](https://matchbox.psdn.io/) (running
    on one of the Raspberry Pis).
  - Matchbox is configured through Terraform.
  - The Kubernetes configuration is all managed by [Flux](https://fluxcd.io/).
