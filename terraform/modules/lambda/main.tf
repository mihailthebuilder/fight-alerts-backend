locals {
  module_prepend = "${var.product}-${var.module}-${var.environment}"
}

resource "aws_s3_bucket" "fight_alerts_scraper_lambda" {
  bucket = "${local.module_prepend}-bucket"

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
  function_name = local.module_prepend

  s3_bucket = aws_s3_bucket.fight_alerts_scraper_lambda.id
  s3_key    = aws_s3_bucket_object.fight_alerts_scraper_lambda.key

  runtime = "go1.x"
  handler = "scraper"
  timeout = 60

  reserved_concurrent_executions = 1

  source_code_hash = data.archive_file.fight_alerts_scraper_lambda.output_base64sha256

  role = aws_iam_role.fight_alerts_scraper_iam_policy.arn
  vpc_config {
    subnet_ids         = [var.subnet_ids[0].public, var.subnet_ids[0].private, var.subnet_ids[1].public, var.subnet_ids[1].private]
    security_group_ids = [aws_security_group.lambda_function_security_group.id]
  }

  environment {
    variables = {
      RDS_HOST     = data.aws_ssm_parameter.db_endpoint.value
      RDS_PASSWORD = var.db_password
      RDS_USERNAME = var.db_username
    }
  }
}

resource "aws_cloudwatch_log_group" "fight_alerts_scraper_lambda" {
  name = "/aws/lambda/${aws_lambda_function.fight_alerts_scraper_lambda.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "fight_alerts_scraper_iam_policy" {
  name = "${local.module_prepend}-iam-role"

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
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.fight_alerts_scraper_iam_policy.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_security_group" "lambda_function_security_group" {
  name   = "${local.module_prepend}-sg"
  vpc_id = var.vpc_id
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group_rule" "allow_inbound" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.lambda_function_security_group.id
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = var.vpc_id
}

resource "aws_route_table" "public" {
  vpc_id = var.vpc_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internet_gateway.id
  }
}

module "lambda_networking_0" {
  public_subnet_id  = var.subnet_ids[0].public
  private_subnet_id = var.subnet_ids[0].private

  source                = "./lambda_networking"
  vpc_id                = var.vpc_id
  public_route_table_id = aws_route_table.public.id
}

module "lambda_networking_1" {
  public_subnet_id  = var.subnet_ids[1].public
  private_subnet_id = var.subnet_ids[1].private

  source                = "./lambda_networking"
  vpc_id                = var.vpc_id
  public_route_table_id = aws_route_table.public.id
}
