################################################################################
# ArgoCD
################################################################################

resource "helm_release" "argocd" {
  name = "argocd"

  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  namespace  = "argocd"
  version    = var.argocd_version

  create_namespace = true
  cleanup_on_fail  = true
  timeout          = 8 * 60 # 8 min

  values = [
    templatefile("${path.module}/templates/argocd-values.yml.tpl", {})
  ]

  depends_on = [
    module.eks
  ]
}

resource "kubectl_manifest" "argocd_config_repository" {
  yaml_body = templatefile("${path.module}/templates/argocd-config-repo.yml.tpl", {
    "username" : data.aws_ssm_parameter.github_username.value,
    "password" : data.aws_ssm_parameter.github_token.value
  })

  depends_on = [
    helm_release.argocd
  ]
}

resource "kubectl_manifest" "argocd_cluster_init" {
  yaml_body = templatefile("${path.module}/templates/cluster-init.yml.tpl", {})

  depends_on = [
    kubectl_manifest.argocd_config_repository
  ]
}
