data "aws_ssm_parameter" "db_endpoint" {
  name = "/${var.product}/${var.environment}/db/endpoint"
}
