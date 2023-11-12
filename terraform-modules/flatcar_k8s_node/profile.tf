resource "matchbox_profile" "install" {
  name = "${var.name}-install-profile"

  kernel = "/assets/flatcar/${var.os_version}/flatcar_production_pxe.vmlinuz"

  initrd = [
    "/assets/flatcar/${var.os_version}/flatcar_production_pxe_image.cpio.gz",
  ]

  args = [
    "initrd=flatcar_production_pxe_image.cpio.gz",
    "flatcar.config.url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "flatcar.first_boot=yes",
    "console=tty0",
    "console=ttyS0",
  ]

  container_linux_config = local.install_config
}

data "ct_config" "node_ignition" {
  content = local.node_config
  strict  = true

  snippets = concat(
    [
      var.role == "controller" || var.role == "primary-controller" ? local.controller_config : "---",
      var.role == "primary-controller" ? local.primary_controller_config : "---",
      var.role == "controller" ? local.secondary_controller_config : "---",
    ],
    var.snippets
  )
}

resource "matchbox_profile" "controller" {
  name         = "${var.name}-controller-profile"
  raw_ignition = data.ct_config.node_ignition.rendered
}
