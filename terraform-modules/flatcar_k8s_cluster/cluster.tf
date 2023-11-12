locals {
  controller_ip_addresses = [for k, v in var.controllers : v.ip_address]
  primary_controller      = var.controllers[var.primary_controller]
  secondary_controllers   = { for k, v in var.controllers : k => v if k != var.primary_controller }
}

module "primary_controller" {
  source                  = "../flatcar_k8s_node"
  name                    = var.primary_controller
  role                    = "primary-controller"
  k8s_api_vip             = var.k8s_api_vip
  k8s_api_address         = var.k8s_api_address
  ssh_authorized_key      = var.ssh_authorized_key
  mac                     = local.primary_controller.mac
  install_disk            = local.primary_controller.install_disk
  secrets_drive           = local.primary_controller.secrets_drive
  snippets                = local.primary_controller.snippets
  controller_ip_addresses = local.controller_ip_addresses
  oidc_issuer_url         = var.oidc_issuer_url
}

module "secondary_controllers" {
  for_each                = local.secondary_controllers
  source                  = "../flatcar_k8s_node"
  name                    = each.key
  role                    = "controller"
  k8s_api_vip             = var.k8s_api_vip
  k8s_api_address         = var.k8s_api_address
  ssh_authorized_key      = var.ssh_authorized_key
  mac                     = each.value.mac
  install_disk            = each.value.install_disk
  secrets_drive           = each.value.secrets_drive
  snippets                = each.value.snippets
  controller_ip_addresses = local.controller_ip_addresses
  oidc_issuer_url         = var.oidc_issuer_url
}

data "external" "cluster_provisioned" {
  program = ["${path.module}/external/wait-for-provisioned"]
  query = {
    primary_controller_ip = local.primary_controller.ip_address
  }
  depends_on = [
    module.primary_controller
  ]
}

resource "null_resource" "provision_secondary_controllers" {
  depends_on = [
    data.external.cluster_provisioned,
  ]

  provisioner "local-exec" {
    command = "${path.module}/external/provision-secondary-controllers ${local.primary_controller.ip_address} ${join(" ", [for k, v in local.secondary_controllers : v.ip_address])}"
  }
}
