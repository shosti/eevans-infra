resource "aws_ses_domain_identity" "domain" {
  domain = var.domain
}

resource "aws_ses_email_identity" "emails" {
  email    = each.key
  for_each = var.emails
}

resource "aws_ses_domain_dkim" "dkim" {
  domain = aws_ses_domain_identity.domain.domain
}

resource "cloudflare_record" "dkim" {
  count      = var.add_dkim_records ? length(aws_ses_domain_dkim.dkim.dkim_tokens) : 0
  zone_id    = var.domain_zone_id
  name       = "${aws_ses_domain_dkim.dkim.dkim_tokens[count.index]}._domainkey"
  type       = "CNAME"
  ttl        = "600"
  value      = "${aws_ses_domain_dkim.dkim.dkim_tokens[count.index]}.dkim.amazonses.com"
  depends_on = [aws_ses_domain_dkim.dkim]
}

resource "aws_ses_domain_mail_from" "mail_domain" {
  domain           = aws_ses_domain_identity.domain.domain
  mail_from_domain = "mail.${aws_ses_domain_identity.domain.domain}"
}

resource "aws_ses_domain_mail_from" "email" {
  domain           = each.value.email
  mail_from_domain = "mail.${aws_ses_domain_identity.domain.domain}"
  for_each         = aws_ses_email_identity.emails
}

resource "cloudflare_record" "mail_from_mx" {
  zone_id  = var.domain_zone_id
  name     = "mail"
  type     = "MX"
  ttl      = "600"
  priority = "10"
  value    = "feedback-smtp.${var.ses_region}.amazonses.com"
}

resource "cloudflare_record" "mail_from_txt" {
  zone_id = var.domain_zone_id
  name    = "mail"
  type    = "TXT"
  ttl     = "600"
  value   = "v=spf1 include:amazonses.com -all"
}
