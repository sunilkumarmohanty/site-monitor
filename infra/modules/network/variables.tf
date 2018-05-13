variable "vpc_cidr" {}
variable "environment" {}

variable "private_subnet_cidrs" {
  type = "list"
}

variable "public_subnet_cidrs" {
  type = "list"
}

variable "availability_zones" {
  type = "list"
}

variable "destination_cidr_block" {
  default = "0.0.0.0/0"
}

variable depends_id {}
