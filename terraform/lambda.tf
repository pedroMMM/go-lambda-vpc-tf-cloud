locals {
  lambda_handler = "lambda"
  lambda_bin     = format("%s/%s", path.module, local.lambda_handler)
  lambda_zip     = "${local.lambda_bin}.zip"
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
      state_s3_key = "state"
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
  name = local.app

  assume_role_policy = data.aws_iam_policy_document.counter.json
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

resource "aws_iam_role_policy_attachment" "counter" {
  role       = aws_iam_role.counter.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
