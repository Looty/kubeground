terraform {
  required_version = ">= 1.3.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.40"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.14.0"
    }
    kubectl = {
      source  = "alekc/kubectl"
      version = "2.0.4"
    }
  }

  backend "s3" {
    bucket  = "kubeground"
    key     = "state/terraform.tfstate"
    region  = "us-east-1"
    profile = "dev"
  }
}

provider "aws" {
  region  = var.region
  profile = "dev"

  default_tags {
    tags = var.tags
  }
}

provider "aws" {
  alias   = "east2"
  region  = var.vpc_region
  profile = "dev"

  default_tags {
    tags = var.tags
  }
}
