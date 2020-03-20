locals {
  bucket_name            = format("%s-%s", local.account_id, local.app)
  bucket_object_wildcard = format("%s/*", aws_s3_bucket.state.arn)
}

resource "aws_s3_bucket" "state" {
  bucket = local.bucket_name
  acl    = "private"
  region = var.region

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.base.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }

  tags = {
    Name  = local.bucket_name
    App   = local.app
    Owner = var.owner
  }
}

resource "aws_s3_bucket_policy" "state" {
  bucket = aws_s3_bucket.state.bucket
  policy = data.aws_iam_policy_document.state.json
}


data "aws_iam_policy_document" "state" {
  statement {
    sid     = "EnforceSecureTransport"
    effect  = "Deny"
    actions = ["s3:*"]

    resources = [local.bucket_object_wildcard]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"

      values = ["false"]
    }
  }

  statement {
    sid     = "DenyIncorrectEncryptionHeader"
    effect  = "Deny"
    actions = ["s3:PutObject"]

    resources = [local.bucket_object_wildcard]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "StringEquals"
      variable = "s3:x-amz-server-side-encryption"

      values = ["AES256"]
    }
  }

  statement {
    sid     = "DenyIncorrectKey"
    effect  = "Deny"
    actions = ["s3:PutObject"]

    resources = [local.bucket_object_wildcard]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "StringNotEqualsIfExists"
      variable = "s3:x-amz-server-side-encryption-aws-kms-key-id"

      values = [aws_kms_key.base.arn]
    }
  }
}
