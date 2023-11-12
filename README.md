# My Personal Infra Repo

This is my infrastructure-as-code repo for my homelab and other self-hosted
stuff. Some more sensitive stuff is in private submodules, so a lot of things
won't work out of the box.

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
  - There are a few VMs running on [KubeVirt](https://kubevirt.io/), including a
    [TrueNAS](https://www.truenas.com/) storage box.

## Interesting Paths

Some places to start:

- `/ansible`: This includes a bunch of Ansible roles and modules. There are some
  roles for [Pi-hole](https://pi-hole.net/) DNS setup (with `keepalived`-based
  failover), [`tangd`](https://github.com/latchset/tang) servers, and
  [VyOS](https://vyos.io) router configuration.
- `/terraform-modules`: This includes some Terraform modules (the top-level
  terraform module is private, unfortunately). The most interesting modules are
  `flatcar_k8s_cluster` and `flatcar_k8s_node`, which declaratively set up
  bare-metal Flatcar Linux Kubernetes worker nodes to be provisioned over
  netboot.
- `/k8s/deploy`: This is where the FluxCD declarative for my homelab Kubernetes
  cluster, which is where almost everything I run ends up going.
