resource "aws_ses_email_identity" "email_identity" {
  email = var.kubeground_inquiries_email
}

resource "aws_iam_policy" "email_policy" {
  name        = "${var.name}AllowSendingEmail"
  description = "Policy to allow sending emails via SES"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ses:SendEmail",
          "ses:SendRawEmail"
        ]
        Effect = "Allow"
        Resource = [
          aws_ses_email_identity.email_identity.arn,
        ]
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "email_policy_attachment" {
  role       = module.eks.eks_managed_node_groups["cluster"].iam_role_name
  policy_arn = aws_iam_policy.email_policy.arn
}
