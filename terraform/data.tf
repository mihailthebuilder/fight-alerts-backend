data "aws_ssm_parameter" "db_username" {
  name = "/${var.product}/${var.environment}/db/username"
}

data "aws_ssm_parameter" "db_password" {
  name = "/${var.product}/${var.environment}/db/password"
}
