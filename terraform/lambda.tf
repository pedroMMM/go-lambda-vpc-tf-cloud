locals {
  lambda_handler   = "lambda"
  lambda_bin       = format("%s/%s", path.module, local.lambda_handler)
  lambda_zip       = "${local.lambda_bin}.zip"
  state_s3_key     = "state"
  state_s3_key_arn = format("%s/%s", aws_s3_bucket.state.arn, local.state_s3_key)
}

resource "aws_lambda_function" "counter" {
  filename         = data.archive_file.counter.output_path
  function_name    = local.app
  role             = aws_iam_role.counter.arn
  handler          = local.lambda_handler
  source_code_hash = filebase64sha256(data.archive_file.counter.output_path)
  runtime          = "go1.x"
  memory_size      = 128

  environment {
    variables = {
      bucket_name  = aws_s3_bucket.state.bucket
      state_s3_key = local.state_s3_key
    }
  }

  tags = {
    Name  = local.app
    App   = local.app
    Owner = var.owner
  }
}

data "archive_file" "counter" {
  type        = "zip"
  source_file = local.lambda_bin
  output_path = local.lambda_zip
}

resource "aws_iam_role" "counter" {
  name               = local.app
  assume_role_policy = data.aws_iam_policy_document.counter.json

  tags = {
    Name  = local.app
    App   = local.app
    Owner = var.owner
  }
}

data "aws_iam_policy_document" "counter" {
  statement {
    sid     = "AllowLambda"
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "counter_basic_exection" {
  role       = aws_iam_role.counter.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "counter_execution" {
  statement {
    sid       = "AllowListVpcs"
    effect    = "Allow"
    actions   = ["ec2:DescribeVpcs"]
    resources = ["*"]
  }

  statement {
    sid       = "AllowListRegions"
    effect    = "Allow"
    actions   = ["ec2:DescribeRegions"]
    resources = ["*"]
  }

  statement {
    sid    = "AllowS3StateUse"
    effect = "Allow"

    actions = [
      "s3:PutObject",
      "s3:GetObject",
    ]

    resources = [local.state_s3_key_arn]
  }

  statement {
    sid    = "AllowS3StateCheckExistance"
    effect = "Allow"

    actions = [
      "s3:ListBucket",
    ]

    resources = [aws_s3_bucket.state.arn]
  }

  statement {
    sid    = "AllowKmsUse"
    effect = "Allow"

    actions = [
      "kms:Decrypt",
      "kms:Encrypt",
    ]

    resources = [aws_kms_key.base.arn]
  }
}

resource "aws_iam_policy" "counter_execution" {
  name   = local.app
  policy = data.aws_iam_policy_document.counter_execution.json
}

resource "aws_iam_role_policy_attachment" "counter_execution" {
  role       = aws_iam_role.counter.name
  policy_arn = aws_iam_policy.counter_execution.arn
}
