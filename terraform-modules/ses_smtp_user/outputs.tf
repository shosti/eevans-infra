output "smtp_user" {
  value = aws_iam_access_key.access_key.id
}

output "smtp_password" {
  value     = data.external.smtp_password.result["key"]
  sensitive = true
}
