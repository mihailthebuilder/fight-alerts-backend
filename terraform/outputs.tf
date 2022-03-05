output "fight_alerts_scraper_s3_for_lambda" {
  description = "Name of the S3 bucket used to store function code."
  value       = aws_s3_bucket.fight_alerts_scraper_lambda.id
}

output "fight_alerts_scraper_lambda" {
  description = "Name of the lambda function."
  value       = aws_lambda_function.fight_alerts_scraper_lambda.function_name
}
