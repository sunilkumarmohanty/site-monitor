variable "name" {}
variable "environment" {}

variable "vpc_id" {}

variable "subnet_cidrs" {
  type = "list"
}

variable "availability_zones" {
  type = "list"
}
