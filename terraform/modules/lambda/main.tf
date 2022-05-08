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

  role = aws_iam_role.fight_alerts_scraper_iam_role.arn
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

resource "aws_iam_role" "fight_alerts_scraper_iam_role" {
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

resource "aws_iam_role_policy_attachment" "execution_policy_attachment" {
  role       = aws_iam_role.fight_alerts_scraper_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy_attachment" "other_policy_attachment" {
  role       = aws_iam_role.fight_alerts_scraper_iam_role.name
  policy_arn = aws_iam_policy.other_policy.arn
}

resource "aws_iam_policy" "other_policy" {
  name   = "${local.module_prepend}-other-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.other_policy_document.json
}

data "aws_iam_policy_document" "other_policy_document" {
  statement {
    effect = "Allow"
    actions = [
      "events:PutRule",
      "events:PutTargets",
      "events:RemoveTargets",
      "events:DeleteRule",
      "events:ListRules",
      "events:ListTargetsByRule"
    ]
    resources = ["*"]
  }
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

module "lambda_networking" {
  count      = 2
  subnet_ids = var.subnet_ids[count.index]

  source                = "./lambda_networking"
  vpc_id                = var.vpc_id
  public_route_table_id = aws_route_table.public.id
}
