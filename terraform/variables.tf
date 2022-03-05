variable "resource_tags" {
  type = object({
    Owner       = string
    Contact     = string
    Environment = string
  })

  default = {
    Owner       = "Mihail_Marian"
    Contact     = "m.marian@elsevier.com"
    Environment = "prod"
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
