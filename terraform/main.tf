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

  tags = var.resource_tags

  acl           = "private"
  force_destroy = true
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


resource "aws_lambda_function" "fight_alerts_scraper_lambda" {
  function_name = "fight-alerts-scraper-lambda-prod"

  s3_bucket = aws_s3_bucket.fight_alerts_scraper_lambda.id
  s3_key    = aws_s3_bucket_object.fight_alerts_scraper_lambda.key

  runtime = "go1.x"
  handler = "scraper"

  source_code_hash = data.archive_file.fight_alerts_scraper_lambda.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
  tags = var.resource_tags
}

resource "aws_cloudwatch_log_group" "fight_alerts_scraper_lambda" {
  name = "/aws/lambda/${aws_lambda_function.fight_alerts_scraper_lambda.function_name}"

  retention_in_days = 30
  tags              = var.resource_tags
}

resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      }
    ]
  })

  tags = var.resource_tags
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
