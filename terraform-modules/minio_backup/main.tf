locals {
  one_day  = 24 * 60 * 60
  one_hour = 60 * 60
}

data "minio_iam_policy_document" "user_policy" {
  statement {
    actions = [
      "s3:ListBucket",
      "s3:DeleteObject",
      "s3:GetObject",
      "s3:PutObject",
    ]
    resources = [
      "arn:aws:s3:::${var.bucket_name}",
      "arn:aws:s3:::${var.bucket_name}/*",
    ]
  }
}

resource "minio_iam_user" "user" {
  name = var.user_name
}

resource "minio_iam_service_account" "service_account" {
  target_user = minio_iam_user.user.name
  lifecycle {
    ignore_changes = [policy]
  }
}

resource "minio_s3_bucket" "bucket" {
  bucket = var.bucket_name
}

resource "minio_iam_policy" "policy" {
  name   = "${var.user_name}-policy"
  policy = data.minio_iam_policy_document.user_policy.json
}

resource "minio_iam_user_policy_attachment" "user_policy" {
  user_name   = minio_iam_user.user.id
  policy_name = minio_iam_policy.policy.id
}

resource "healthchecksio_check" "backup" {
  name    = "${var.bucket_name}-backup"
  timeout = local.one_day
  grace   = local.one_hour
}
