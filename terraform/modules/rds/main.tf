locals {
  module_prepend = "${var.product}-${var.module}-${var.environment}"
  ingress_ips    = jsondecode(data.aws_ssm_parameter.ingress_ips.value)
}

resource "aws_security_group" "rds_cluster_security_group" {
  name   = "${local.module_prepend}-cluster-security-group"
  vpc_id = var.vpc_id
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_subnet_group" "rds_db_subnet_group" {
  name       = "${local.module_prepend}-subnet-group"
  subnet_ids = [var.subnet_ids[0].public, var.subnet_ids[1].public]
}

resource "aws_rds_cluster" "rds_cluster" {
  apply_immediately    = true
  cluster_identifier   = "${local.module_prepend}-cluster"
  engine               = "aurora-postgresql"
  master_password      = var.db_password
  master_username      = var.db_username
  db_subnet_group_name = aws_db_subnet_group.rds_db_subnet_group.name

  # below lines are used to ensure Terraform can destroy RDS instance successfully when needed
  final_snapshot_identifier = "dummy"
  skip_final_snapshot       = true

  vpc_security_group_ids = [aws_security_group.rds_cluster_security_group.id]
}

resource "aws_security_group_rule" "allow_inbound" {
  count             = length(local.ingress_ips)
  description       = element(keys(local.ingress_ips), count.index)
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  cidr_blocks       = split(",", lookup(local.ingress_ips, element(keys(local.ingress_ips), count.index)))
  security_group_id = aws_security_group.rds_cluster_security_group.id
}

resource "aws_rds_cluster_instance" "rds_instance" {
  apply_immediately            = true
  identifier                   = "${local.module_prepend}-cluster-instance"
  cluster_identifier           = aws_rds_cluster.rds_cluster.cluster_identifier
  instance_class               = "db.t3.medium"
  engine                       = "aurora-postgresql"
  publicly_accessible          = true
  performance_insights_enabled = true
  depends_on = [
    aws_rds_cluster.rds_cluster
  ]
}
