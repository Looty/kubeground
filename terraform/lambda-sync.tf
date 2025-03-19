// ################################################################################
// # AWS Lambda sync (Create/delete users when mentor uploads .csv of students)
// ################################################################################

resource "aws_s3_bucket_notification" "approved_emails" {
  bucket = data.aws_s3_bucket.bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.sync_lambda.arn
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = "${local.csv_sync_file}"
  }
}

resource "aws_iam_role" "sync_lambda_role" {
  name = "${var.name}_sync_lambda_role"

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

resource "aws_iam_role_policy" "sync_lambda_policy" {
  name = "${var.name}_sync_lambda_policy"
  role = aws_iam_role.sync_lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ],
        Effect = "Allow",
        Resource = [
          data.aws_s3_bucket.bucket.arn,
          "${data.aws_s3_bucket.bucket.arn}/*"
        ]
      },
      {
        Action = [
          "cognito-idp:AdminCreateUser",
          "cognito-idp:AdminDeleteUser",
          "cognito-idp:AdminSetUserPassword",
          "cognito-idp:ListUsers"
        ],
        Effect   = "Allow",
        Resource = "*"
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

data "archive_file" "sync_lambda_code" {
  type        = "zip"
  source_dir  = "../lambda/sync/code/"
  output_path = "../lambda/output/sync_lambda_function_payload.zip"
}

resource "aws_lambda_function" "sync_lambda" {
  function_name    = "${var.name}-sync"
  filename         = data.archive_file.sync_lambda_code.output_path
  source_code_hash = data.archive_file.sync_lambda_code.output_base64sha256
  role             = aws_iam_role.sync_lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  publish          = false
  timeout          = 900 # set maximum - depends on amount of users

  runtime      = "python3.12"
  memory_size  = 128
  package_type = "Zip"
  skip_destroy = false

  ephemeral_storage { size = 512 }
  architectures = ["x86_64"]

  environment {
    variables = {
      USER_POOL_ID     = aws_cognito_user_pool.email_user_pool.id
      DEFAULT_PASSWORD = data.aws_ssm_parameter.cognito_default_password.value
      CSV_PATH         = "customers/${local.csv_sync_file}"
    }
  }

  logging_config {
    log_format = "Text"
    log_group  = aws_cloudwatch_log_group.sync_lambda_group.name
  }

  depends_on = [
    aws_cloudwatch_log_group.sync_lambda_group,
  ]
}

resource "aws_lambda_permission" "allow_s3" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.sync_lambda.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = data.aws_s3_bucket.bucket.arn
}

resource "aws_cloudwatch_log_group" "sync_lambda_group" {
  name              = "/aws/lambda/${var.name}-sync"
  retention_in_days = 0
}
