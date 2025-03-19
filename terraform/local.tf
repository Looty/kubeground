locals {
  external_dns_credentials = <<-EOT
[default]
aws_access_key_id = ${aws_iam_access_key.external_dns.id}
aws_secret_access_key = ${aws_iam_access_key.external_dns.secret}
EOT

  cognito_auth_url          = "https://${var.name}.auth.${var.region}.amazoncognito.com/login?response_type=code&client_id=${aws_cognito_user_pool_client.email_user_pool_client.id}&redirect_uri=${tolist(aws_cognito_user_pool_client.email_user_pool_client.callback_urls)[1]}"
  kubeground_login_page_url = "login.${var.kubeground_domain_dns}"
  csv_sync_file             = "approved_emails.csv"
  login_bucket_name         = "${var.name}-login"
  login_bucket_url          = "https://${local.login_bucket_name}.s3.${var.region}.amazonaws.com/${var.kubeground_login_site_index}"
}
