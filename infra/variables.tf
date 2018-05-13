variable "region" {
  default = "eu-west-1"
}

variable "cluster_name" {}
variable "environment" {}
variable "vpc_cidr" {}

variable "public_subnet_cidrs" {
  type = "list"
}

variable "private_subnet_cidrs" {
  type = "list"
}

variable "availability_zones" {
  type = "list"
}

variable "efs_task_mount_path" {}

variable "efs_host_mount_path" {}
variable "log_retention_days" {}

variable "monitor_config_file_path" {}
variable "monitor_user_agent" {}
variable "monitor_timeout" {}

variable "aws_ami" {}

variable "efs_name" {}
variable "monitor_task_name" {}

variable "task_image" {}
variable "task_image_version" {}
