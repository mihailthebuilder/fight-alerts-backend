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

resource "aws_s3_bucket" "fight_alerts_scraper_lambda" {
  bucket = "fight-alerts-scraper-lambda-bucket-prod"

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
  value       = aws_s3_bucket.fight_alerts_scraper_lambda.id
}

data "archive_file" "fight_alerts_scraper_lambda" {
  type        = "zip"
  source_file = "../functions/bin/scraper"
  output_path = "bin/scraper.zip"
}

resource "aws_s3_bucket_object" "fight_alerts_scraper_lambda" {
  bucket = aws_s3_bucket.fight_alerts_scraper_lambda.id

  key    = "scraper.zip"
  source = data.archive_file.fight_alerts_scraper_lambda.output_path

  etag = filemd5(data.archive_file.fight_alerts_scraper_lambda.output_path)
}