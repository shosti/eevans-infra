variable "name" {
  type = string
}

variable "mac" {
  type = string
}

variable "k8s_api_vip" {
  type        = string
  description = "Kubernetes API virtual IP address"
}

variable "controller_ip_addresses" {
  type = list(string)
}

variable "k8s_api_address" {
  type        = string
  description = "Kubernetes API FQDN (should resolve to the VIP)"
}

variable "role" {
  type = string
  validation {
    condition = anytrue([
      var.role == "primary-controller",
      var.role == "controller",
      # var.role == "worker",
    ])
    error_message = "Node role is required."
  }
}

variable "channel" {
  type    = string
  default = "stable"
}

variable "os_version" {
  type    = string
  default = "3602.2.0"
}

variable "k8s_version" {
  type    = string
  default = "v1.29.0"
}

variable "matchbox_http_endpoint" {
  type    = string
  default = "https://matchbox.eevans.me:8880"
}

variable "ssh_authorized_key" {
  type = string
}

variable "install_disk" {
  type    = string
  default = "/dev/sda"
}

variable "secrets_drive" {
  type        = string
  description = "USB thumb drive on which to store secrets (blank to disable). Should be in /dev/disk/by-id path."
}

variable "pod_subnet" {
  type    = string
  default = "10.101.0.0/16"
}

variable "service_subnet" {
  type    = string
  default = "10.70.0.0/16"
}

variable "oidc_issuer_url" {
  type = string
}

variable "snippets" {
  type        = list(string)
  description = "Extra ct snippets for the node's configuration"
  default     = []
}

variable "encrypted_directories" {
  type        = list(string)
  description = "Directories to encrypt on the host"
  default = [
    "/opt/etc/kubernetes",
    "/opt/local-path-provisioner",
    "/opt/secrets",
    "/var/lib/etcd",
    "/var/lib/kubelet/pki",
    "/var/lib/rook",
  ]
}
