resource "aws_cloudwatch_event_rule" "cron" {
  name                = local.app
  description         = format("Lambda CRON trigger for %s", local.app)
  schedule_expression = format("cron(%s)", var.cron)
  is_enabled          = true

  tags = {
    Name  = local.app
    App   = local.app
    Owner = var.owner
  }
}

resource "aws_cloudwatch_event_target" "lambda_trigger" {
  rule      = aws_cloudwatch_event_rule.cron.name
  target_id = local.app
  arn       = aws_lambda_function.counter.arn
}

resource "aws_lambda_permission" "lambda_trigger" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.counter.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cron.arn
}
