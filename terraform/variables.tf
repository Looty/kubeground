variable "region" {
  description = "The AWS region to deploy into (e.g. us-east-1)"
  type        = string
  default     = "us-east-1"
}

variable "vpc_region" {
  description = "The AWS region to deploy into (e.g. us-east-1)"
  type        = string
  default     = "us-east-2"
}

variable "name" {
  description = "Project name"
  type        = string
  default     = "kubeground"
}

variable "cluster_version" {
  description = "EKS version"
  type        = string
  default     = "1.29" # latest = 1.30
}

variable "cluster_node_type" {
  description = "EKS instance node type"
  type        = string
  default     = "t3.medium"
}

variable "cluster_node_limits" {
  description = "Map of minimum and maximum number of nodes in the cluster"
  type        = map(number)
  default = {
    min          = 3
    desired_size = 3
    max          = 6
  }
}

variable "tags" {
  description = "Terraform tags"
  type        = map(any)
  default = {
    Environment = "kubeground"
    GithubRepo  = "https://github.com/Looty/kubeground"
    Terraform   = "true"
    project     = "k8s_kubeground"
    startdate   = "01/06/24"
    enddate     = "N/A"
    Expiration  = "N/A"
    Objective   = "Internal project"
    Email       = "erez.mizrahi@gmail.com"
  }
}

variable "argocd_version" {
  description = "ArgoCD version"
  type        = string
  default     = "7.1.4"
}

variable "externaldns_namespace" {
  description = "ExternalDNS namespace name"
  type        = string
  default     = "external-dns"
}

variable "kubeground_inquiries_email" {
  description = "Email address to be sent from & help customers"
  type        = string
  default     = "erez.mizrahi@gmail.com"
}

variable "kubeground_domain_dns" {
  description = "Kubeground domain DNS"
  type        = string
  default     = "kubegroundapp.click"
}

variable "kubeground_login_site_index" {
  description = "Kubeground S3 login site index"
  type        = string
  default     = "index.html"
}

variable "kubeground_virtual_platform_namespace" {
  description = "Kubeground virtual platform namespace"
  type        = string
  default     = "virtual-platform"
}

variable "lambda_cleaner_trigger" {
  description = "The scheduling expression in minutes. For example, cron(0 20 * * ? *) or rate(5 minutes). At least one of schedule_expression or event_pattern is required. Can only be used on the default event bus."
  type        = string
  default     = "60"
}

variable "lambda_cleaner_delete_customer_after" {
  description = "Delete old customer configs after x minutes"
  type        = string
  default     = "300"
}
