locals {
  kms_access_roles = [aws_iam_role.counter.arn]
  kms_alias        = format("alias/%s", local.app)
}

resource "aws_kms_key" "base" {
  description             = "Just a base KMS for ${local.app} App"
  deletion_window_in_days = 30
  policy                  = data.aws_iam_policy_document.base.json
  enable_key_rotation     = true

  tags = {
    Name  = local.app
    App   = local.app
    Owner = var.owner
  }
}

resource "aws_kms_alias" "base" {
  name          = local.kms_alias
  target_key_id = aws_kms_key.base.arn
}

data "aws_iam_policy_document" "base" {
  statement {
    sid    = "EnableIamUserPermissions"
    effect = "Allow"

    actions = [
      "kms:*",
    ]

    resources = [
      "*",
    ]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${local.account_id}:root"]
    }
  }

  statement {
    sid    = "AllowUseOfTheKey"
    effect = "Allow"

    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey",
      "kms:CreateGrant",
      "kms:ListGrants",
      "kms:RevokeGrant",
    ]

    resources = [
      "*",
    ]

    principals {
      type = "AWS"

      identifiers = local.kms_access_roles
    }
  }
}

