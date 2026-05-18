variable "aws_region_virginia" {
  type = string
}

variable "virginia_vpc" {
  type = object({
    name  = string
    env   = string
    owner = string
    cidr  = string
  })
}

variable "subnets" {
  description = "List Subnets"
  type        = list(string)
}

variable "tags" {
  description = "Tags of Project"
  type        = map(string)
}
