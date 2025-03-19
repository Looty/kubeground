apiVersion: v1
kind: Secret
metadata:
  name: config-repo
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  url: https://github.com/Looty/kubeground-config
  project: default
  type: git
  username: ${username}
  password: ${password}
