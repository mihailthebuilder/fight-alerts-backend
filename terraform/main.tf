terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.37"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "fight_alerts_scraper_lambda_bucket" {
  bucket = "fight_alerts_scraper_lambda_bucket-prod"

  tags = {
    Owner       = "Mihail_Marian"
    Contact     = "m.marian@elsevier.com"
    Environment = "prod"
  }

  acl           = "private"
  force_destroy = true
}

output "fight_alerts_scraper_lambda_function" {
  description = "Name of the S3 bucket used to store function code."
  value       = aws_s3_bucket.fight_alerts_scraper_lambda_bucket.id
}