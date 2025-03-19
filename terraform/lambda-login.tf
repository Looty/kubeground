################################################################################
# AWS Lambda Login
################################################################################

resource "aws_iam_role" "login_lambda_role" {
  name = "${var.name}_login_lambda_role"

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

resource "aws_iam_role_policy" "login_lambda_policy" {
  name = "${var.name}_login_lambda_policy"
  role = aws_iam_role.login_lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "lambda:InvokeFunction",
        ],
        Effect   = "Allow",
        Resource = aws_lambda_function.login_lambda.arn
      },
      {
        Action = [
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

data "archive_file" "login_lambda_layer" {
  type        = "zip"
  source_dir  = "../lambda/login/layer/"
  output_path = "../lambda/output/login_lambda_layer_payload.zip"
}

# for GitPython
resource "aws_lambda_layer_version" "login_gitpython_layer" {
  filename         = data.archive_file.login_lambda_layer.output_path
  layer_name       = "python-git"
  description      = "for kubeground"
  source_code_hash = data.archive_file.login_lambda_layer.output_base64sha256

  compatible_runtimes      = ["python3.12"]
  compatible_architectures = ["x86_64"]
}

data "archive_file" "login_lambda_code" {
  type        = "zip"
  source_dir  = "../lambda/login/code/"
  output_path = "../lambda/output/login_lambda_function_payload.zip"
}

resource "aws_lambda_function" "login_lambda" {
  function_name    = "${var.name}-login"
  filename         = data.archive_file.login_lambda_code.output_path
  source_code_hash = data.archive_file.login_lambda_code.output_base64sha256
  role             = aws_iam_role.login_lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  publish          = false
  timeout          = 20
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
      GIT_PYTHON_REFRESH = "quiet"
      GITHUB_TOKEN       = data.aws_ssm_parameter.github_token.value
    }
  }

  logging_config {
    log_format = "Text"
    log_group  = aws_cloudwatch_log_group.login_lambda_group.name
  }

  depends_on = [
    aws_cloudwatch_log_group.login_lambda_group,
  ]
}

resource "aws_cloudwatch_log_group" "login_lambda_group" {
  name              = "/aws/lambda/${var.name}-login"
  retention_in_days = 0
}

resource "kubectl_manifest" "virtual_platform_namespace" {
  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: ${var.kubeground_virtual_platform_namespace}
YAML

  depends_on = [
    module.eks
  ]
}
