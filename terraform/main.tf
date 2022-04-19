terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.37"
    }
  }

  backend "s3" {
    bucket = "fight-alerts-backend-terraform"
    key    = "backend"
    region = "us-east-1"
  }
}

provider "aws" {
  region  = var.region
  profile = "default"

  default_tags {
    tags = var.resource_tags
  }
}

locals {
  db_username = data.aws_ssm_parameter.db_username.value
  db_password = data.aws_ssm_parameter.db_password.value
}

module "lambda" {
  source      = "./modules/lambda"
  environment = var.environment
  region      = var.region
  product     = var.product
  module      = "scraper-lambda"
  subnet_ids  = var.subnet_ids
  vpc_id      = var.vpc_id

  db_username = local.db_username
  db_password = local.db_password
}

module "rds" {
  source      = "./modules/rds"
  environment = var.environment
  product     = var.product
  module      = "rds"
  subnet_ids  = var.subnet_ids
  vpc_id      = var.vpc_id

  db_username = local.db_username
  db_password = local.db_password
}
