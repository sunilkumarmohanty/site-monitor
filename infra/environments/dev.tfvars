region = "eu-west-1"

environment = "dev"

vpc_cidr = "10.0.0.0/16"

public_subnet_cidrs = ["10.0.2.0/24", "10.0.3.0/24"]

private_subnet_cidrs = ["10.0.0.0/24", "10.0.1.0/24"]

availability_zones = ["eu-west-1a", "eu-west-1b"]

cluster_name = "monitor"

monitor_task_name = "monitor"

efs_task_mount_path = "/efs"

efs_host_mount_path = "/efs"

log_retention_days = 5

monitor_config_file_path = "./monitor.config"

monitor_user_agent = "monitorv1.0.0"

monitor_timeout = 30

aws_ami = "ami-2d386654"

efs_name = "monitor_logs"

task_image = "048989072918.dkr.ecr.eu-west-1.amazonaws.com/monitor"

task_image_version = "latest"
