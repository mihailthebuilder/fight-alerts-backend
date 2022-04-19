variable "resource_tags" {
  default = {
    Name        = "fight_alerts_resource"
    Owner       = "Mihail_Marian"
    Contact     = "m.marian@elsevier.com"
    Environment = "prod"
    Product     = "Fight_Alerts"
  }
}

variable "product" {
  type    = string
  default = "fight-alerts"
}

variable "environment" {
  type    = string
  default = "prod"
}

variable "module" {
  type    = string
  default = "lambda-function"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "subnet_ids" {
  default = [
    {
      "public" : "subnet-002681aa4abc15084",
      "private" : "subnet-01f5ac87e00c67cc7"
    },
    {
      "public" : "subnet-05c938c64ff3b334f",
      "private" : "subnet-0da22e85ac7387dfc"
    },
  ]
}

variable "vpc_id" {
  default = "vpc-086415cfebaa3d6d9"
}
