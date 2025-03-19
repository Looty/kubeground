resource "aws_iam_user" "external_dns" {
  name = "${var.name}-externaldns"
}

resource "aws_iam_policy" "external_dns" {
  name        = "${var.name}AllowExternalDNSUpdates"
  description = "Policy for external-dns to manage Route 53 records"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "route53:ChangeResourceRecordSets"
        ]
        Resource = [
          "arn:aws:route53:::hostedzone/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "route53:ListHostedZones",
          "route53:ListResourceRecordSets",
          "route53:ListTagsForResource"
        ]
        Resource = [
          "*"
        ]
      }
    ]
  })
}

resource "aws_iam_user_policy_attachment" "external_dns" {
  user       = aws_iam_user.external_dns.name
  policy_arn = aws_iam_policy.external_dns.arn
}

resource "aws_iam_access_key" "external_dns" {
  user = aws_iam_user.external_dns.name
}

resource "aws_ssm_parameter" "external_dns_access_key" {
  name  = "/${var.name}/externaldns/credentials/accesskey"
  type  = "String"
  value = aws_iam_access_key.external_dns.id
}

resource "aws_ssm_parameter" "external_dns_secret_key" {
  name  = "/${var.name}/externaldns/credentials/secretkey"
  type  = "String"
  value = aws_iam_access_key.external_dns.secret
}

resource "kubectl_manifest" "external_dns_namespace" {
  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: ${var.externaldns_namespace}
YAML

  depends_on = [
    module.eks
  ]
}

resource "kubectl_manifest" "external_dns_secret" {
  yaml_body = <<YAML
apiVersion: v1
kind: Secret
metadata:
  name: ${var.externaldns_namespace}
  namespace: ${var.externaldns_namespace}
type: Opaque
data:
  credentials: ${base64encode(local.external_dns_credentials)}
YAML

  depends_on = [
    module.eks,
    kubectl_manifest.external_dns_namespace
  ]
}
