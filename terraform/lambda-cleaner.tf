################################################################################
# AWS Lambda Cleaner (delete customers older than X)
################################################################################

resource "aws_iam_role" "cleaner_lambda_role" {
  name = "${var.name}_cleaner_lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy" "cleaner_lambda_policy" {
  name = "${var.name}_cleaner_lambda_policy"
  role = aws_iam_role.cleaner_lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "lambda:InvokeFunction",
        ],
        Effect   = "Allow",
        Resource = aws_lambda_function.cleaner_lambda.arn
      },
      {
        Action = [
          "logs:*",
        ]
        Effect   = "Allow",
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

data "archive_file" "cleaner_lambda_code" {
  type        = "zip"
  source_dir  = "../lambda/cleaner/code/"
  output_path = "../lambda/output/cleaner_lambda_function_payload.zip"
}

resource "aws_lambda_function" "cleaner_lambda" {
  function_name    = "${var.name}-cleaner"
  filename         = data.archive_file.cleaner_lambda_code.output_path
  source_code_hash = data.archive_file.cleaner_lambda_code.output_base64sha256
  role             = aws_iam_role.cleaner_lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  publish          = false
  timeout          = 60
  layers = [
    # for GitPython + Git
    aws_lambda_layer_version.login_gitpython_layer.arn,
    "arn:aws:lambda:us-east-1:553035198032:layer:git-lambda2:8",
  ]

  runtime      = "python3.12"
  memory_size  = 128
  package_type = "Zip"
  skip_destroy = false

  ephemeral_storage { size = 512 }
  architectures = ["x86_64"]

  environment {
    variables = {
      TRIGGER_FORMAT        = "minutes"
      TRIGGER_TIME          = var.lambda_cleaner_trigger
      DELETE_CUSTOMER_AFTER = var.lambda_cleaner_delete_customer_after
      GITHUB_TOKEN          = data.aws_ssm_parameter.github_token.value
    }
  }

  logging_config {
    log_format = "Text"
    log_group  = aws_cloudwatch_log_group.cleaner_lambda_group.name
  }

  depends_on = [
    aws_cloudwatch_log_group.cleaner_lambda_group,
  ]
}

resource "aws_cloudwatch_log_group" "cleaner_lambda_group" {
  name              = "/aws/lambda/${var.name}-cleaner"
  retention_in_days = 0
}

resource "aws_cloudwatch_event_rule" "cleaner_event_rule" {
  name                = "${var.name}-cleaner-trigger"
  schedule_expression = "rate(${var.lambda_cleaner_trigger} minutes)"
}

resource "aws_cloudwatch_event_target" "trigger_lambda" {
  rule      = aws_cloudwatch_event_rule.cleaner_event_rule.name
  target_id = "lambda_target"
  arn       = aws_lambda_function.cleaner_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_invoke" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cleaner_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cleaner_event_rule.arn
}
