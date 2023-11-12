locals {
  install_config = templatefile("${path.module}/cl/install.yaml",
    {
      os_channel         = var.channel
      os_version         = var.os_version
      ignition_endpoint  = format("%s/ignition", var.matchbox_http_endpoint)
      install_disk       = var.install_disk
      ssh_authorized_key = var.ssh_authorized_key
      # profile uses -b baseurl to install from matchbox cache
      baseurl_flag = "-b ${var.matchbox_http_endpoint}/assets/flatcar"
    }
  )

  node_config = templatefile("${path.module}/cl/node.yaml",
    {
      hostname           = var.name
      ssh_authorized_key = var.ssh_authorized_key
      k8s_version        = var.k8s_version
      mac                = var.mac
    }
  )

  controller_config = templatefile("${path.module}/cl/controller.yaml",
    {
      hostname              = var.name
      secrets_drive         = var.secrets_drive
      k8s_api_vip           = var.k8s_api_vip
      encrypted_directories = join(" ", var.encrypted_directories)
    }
  )

  primary_controller_config = templatefile("${path.module}/cl/primary-controller.yaml",
    {
      k8s_api_address = var.k8s_api_address
      k8s_version     = var.k8s_version
      pod_subnet      = var.pod_subnet
      service_subnet  = var.service_subnet
      oidc_issuer_url = var.oidc_issuer_url
    }
  )

  secondary_controller_config = templatefile("${path.module}/cl/secondary-controller.yaml", {})
}
