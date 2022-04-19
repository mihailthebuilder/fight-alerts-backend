data "aws_ssm_parameter" "ingress_ips" {
  name = "/${var.product}/${var.environment}/db/ingress-ips"
}

