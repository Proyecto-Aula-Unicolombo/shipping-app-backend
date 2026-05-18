terraform {
  required_providers {
    aws = {
        source = "hashicorp/aws"
        version = "~>4.36.0"
    }
  }
  required_version = "~>1.15.0"
}

provider "aws" {
  region = var.aws_region_virginia
  default_tags {
    tags = var.tags
  }
}

