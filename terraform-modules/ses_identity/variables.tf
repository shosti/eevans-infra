variable "emails" {
  type = set(string)
}

variable "domain" {
  type = string
}

variable "domain_zone_id" {
  type = string
}

variable "ses_region" {
  type    = string
  default = "us-west-2"
}

variable "add_dkim_records" {
  type        = bool
  default     = false
  description = "Terraform is dumb, so you have to run this once set to false and then set it to true"
}
