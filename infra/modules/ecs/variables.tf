variable "environment" {}
variable "cluster_name" {}
variable "depends_id" {}
variable "vpc_id" {}

variable "min_size" {
  default = 1
}

variable "max_size" {
  default = 1
}

variable "desired_capacity" {
  default = 1
}

variable "private_subnet_ids" {
  type = "list"
}

variable "aws_ami" {}
variable "efs_mount_path" {}
variable "region" {}

variable "instance_type" {
  default = "t2.micro"
}

variable "efs_id" {}
