data "aws_iam_policy_document" "send_email" {
  statement {
    effect    = "Allow"
    actions   = ["ses:SendRawEmail"]
    resources = ["*"]
  }
}

resource "aws_iam_user" "user" {
  name = var.username
}

resource "aws_iam_policy" "send_email" {
  name        = "${var.username}-send-email"
  description = "Policy for sending email through SES"
  policy      = data.aws_iam_policy_document.send_email.json
}

resource "aws_iam_user_policy_attachment" "send_email" {
  user       = aws_iam_user.user.name
  policy_arn = aws_iam_policy.send_email.arn
}

resource "aws_iam_access_key" "access_key" {
  user = aws_iam_user.user.name
}

data "external" "smtp_password" {
  program = ["python", "${path.module}/external/aws-smtp-credential.py"]

  query = {
    region            = var.ses_region
    secret_access_key = aws_iam_access_key.access_key.secret
  }
}
