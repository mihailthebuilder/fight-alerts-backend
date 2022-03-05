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
