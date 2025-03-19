################################################################################
# Data
################################################################################

data "aws_s3_bucket" "bucket" {
  bucket = var.name
}

################################################################################
# AWS Parameter Store
################################################################################

data "aws_ssm_parameter" "github_token" {
  name = "/${var.name}/github/token"
}

data "aws_ssm_parameter" "github_username" {
  name = "/${var.name}/github/username"
}

data "aws_ssm_parameter" "cognito_default_password" {
  name = "/${var.name}/cognito/default_password"
}
