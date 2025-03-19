apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: cluster-init
  namespace: argocd
spec:
  goTemplate: true
  goTemplateOptions: ["missingkey=error"]
  generators:
  - list:
      elements:
      - application: crossplane
        path: crossplane/
        namespace: crossplane-system
      - application: virtualplatform
        path: virtualplatforms/
        namespace: virtual-platform
      - application: argo-events
        path: argo-events/
        namespace: argo-events
      - application: ingress-nginx
        path: nginx-controller/
        namespace: ingress-nginx
      - application: external-dns
        path: external-dns/
        namespace: external-dns
      - application: metrics-server
        path: metrics-server/
        namespace: kube-system
  template:
    metadata:
      name: "{{.application}}-config"
    spec:
      project: default
      source:
        repoURL: https://github.com/Looty/kubeground-config.git
        path: "{{.path}}"
        targetRevision: main
      destination:
        server: https://kubernetes.default.svc
        namespace: "{{.namespace}}"
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
        syncOptions:
          - CreateNamespace=true
