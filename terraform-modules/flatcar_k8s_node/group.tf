resource "matchbox_group" "install" {
  name    = "${var.name}-install-group"
  profile = matchbox_profile.install.name

  selector = {
    mac = var.mac
  }
}

resource "matchbox_group" "controller" {
  name    = "${var.name}-controller-group"
  profile = matchbox_profile.controller.name

  selector = {
    mac = var.mac
    os  = "installed"
  }
}
