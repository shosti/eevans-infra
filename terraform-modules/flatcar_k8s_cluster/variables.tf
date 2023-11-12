variable "k8s_api_vip" {
  type        = string
  description = "Kubernetes API virtual IP address"
}

variable "k8s_api_address" {
  type        = string
  description = "DNS address of K8s API (should resolve to the VIP)"
}

variable "ssh_authorized_key" {
  type = string
}

variable "controllers" {
  type = map(object({
    mac           = string
    ip_address    = string
    install_disk  = string
    secrets_drive = string
    snippets      = list(string)
  }))
}

variable "oidc_issuer_url" {
  type = string
}

variable "primary_controller" {
  type        = string
  description = "Node name for the primary controller (which will get provisioned first and generate secrets)."
}
