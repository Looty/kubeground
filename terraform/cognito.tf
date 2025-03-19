################################################################################
# AWS Cognito
################################################################################

resource "aws_iam_policy" "cognito_policy" {
  name = "Cognito-authenticated-1718529745191"
  path = "/service-role/"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "cognito-identity:GetCredentialsForIdentity",
        ]
        Effect = "Allow"
        Resource = [
          "*",
        ]
      },
    ]
  })
}

resource "aws_iam_role" "iam_for_cognito" {
  name = "kubeground"
  path = "/service-role/"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sts:AssumeRoleWithWebIdentity",
        ]
        Effect = "Allow",
        Principal = {
          Federated = "cognito-identity.amazonaws.com"
        }
        Condition = {
          StringEquals = {
            "cognito-identity.amazonaws.com:aud" = aws_cognito_identity_pool.identity_pool.id
          }
          "ForAnyValue:StringLike" = {
            "cognito-identity.amazonaws.com:amr" = "authenticated"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "cognito_policy_attachment" {
  role       = aws_iam_role.iam_for_cognito.name
  policy_arn = aws_iam_policy.cognito_policy.arn
}

resource "aws_cognito_user_pool" "email_user_pool" {
  name                     = "${var.name}-email-user-pool"
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  admin_create_user_config {
    allow_admin_create_user_only = true
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "admin_only"
      priority = 1
    }
  }

  schema {
    name                = "email"
    attribute_data_type = "String"
    mutable             = true
    required            = true
  }

  password_policy {
    minimum_length                   = 6
    require_lowercase                = true
    temporary_password_validity_days = 7
  }

  lambda_config {
    post_authentication = aws_lambda_function.login_lambda.arn
  }
}

resource "aws_cognito_user_pool_client" "email_user_pool_client" {
  name         = "${var.name}-email-user-pool-client"
  user_pool_id = aws_cognito_user_pool.email_user_pool.id

  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["email", "openid", "profile"]
  callback_urls = [
    local.login_bucket_url,
    "http://localhost:3000/api/auth/callback/cognito"
  ]
}

resource "aws_cognito_identity_pool" "identity_pool" {
  identity_pool_name               = var.name
  allow_unauthenticated_identities = false
  allow_classic_flow               = false

  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.email_user_pool_client.id
    provider_name           = aws_cognito_user_pool.email_user_pool.endpoint
    server_side_token_check = false
  }
}

resource "aws_cognito_user_pool_domain" "cognito_domain" {
  domain       = var.name
  user_pool_id = aws_cognito_user_pool.email_user_pool.id
}

################################################################################
# Login bucket
################################################################################

resource "aws_s3_bucket" "login_site" {
  bucket = local.login_bucket_name

  website {
    index_document = var.kubeground_login_site_index
  }
}

resource "aws_s3_bucket_ownership_controls" "login_site_acl_ownership" {
  bucket = aws_s3_bucket.login_site.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }

  depends_on = [
    aws_s3_bucket_public_access_block.login_site_public_access_block
  ]
}

resource "aws_s3_bucket_acl" "login_site_acl" {
  bucket = aws_s3_bucket.login_site.id
  acl    = "public-read"

  depends_on = [
    aws_s3_bucket_ownership_controls.login_site_acl_ownership
  ]
}

resource "aws_s3_bucket_object" "index_html" {
  bucket       = aws_s3_bucket.login_site.bucket
  key          = var.kubeground_login_site_index
  acl          = "public-read"
  content_type = "text/html"

  content = <<-HTML
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>${var.name}</title>
    </head>
    <body>
        <h1>You have successfully logged in!</h1>
        <p>We are working on sending you an email invitation, please stand by.</p>
        <p>This process typically takes about 7 minutes.</p>
    </body>
    </html>
  HTML
}

resource "aws_s3_bucket_public_access_block" "login_site_public_access_block" {
  bucket                  = aws_s3_bucket.login_site.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
