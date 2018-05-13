provider "aws" {
  region  = "${var.region}"
  version = "1.18"
}

module "network" {
  source               = "modules/network"
  environment          = "${var.environment}"
  vpc_cidr             = "${var.vpc_cidr}"
  public_subnet_cidrs  = "${var.public_subnet_cidrs}"
  private_subnet_cidrs = "${var.private_subnet_cidrs}"
  availability_zones   = "${var.availability_zones}"

  # terraform does not currently support depends_on on module. The below is a hack from https://github.com/hashicorp/terraform/issues/1178#issuecomment-207369534
  # this hack ensures that network is created before ecs instances
  depends_id = ""
}

module "efs" {
  source = "modules/efs"

  environment = "${var.environment}"
  name        = "${var.efs_name}"
  subnet_ids  = "${module.network.private_subnet_ids}"
  vpc_id      = "${module.network.vpc_id}"
  count       = "${length(var.private_subnet_cidrs)}"
}

module "cluster" {
  source = "modules/ecs"

  environment  = "${var.environment}"
  region       = "${var.region}"
  cluster_name = "${var.cluster_name}"

  depends_id = "${module.network.depends_id}"

  vpc_id = "${module.network.vpc_id}"

  private_subnet_ids = "${module.network.private_subnet_ids}"
  aws_ami            = "${var.aws_ami}"
  efs_mount_path     = "${var.efs_host_mount_path}"
  efs_id             = "${module.efs.efs_id}"
}

//Create services and tasks

module "tasks" {
  source = "modules/tasks"

  environment              = "${var.environment}"
  region                   = "${var.region}"
  name                     = "${var.monitor_task_name}"
  cluster_id               = "${module.cluster.ecs_cluster_id}"
  vpc_id                   = "${module.network.vpc_id}"
  efs_task_mount_path      = "${var.efs_task_mount_path}"
  efs_host_mount_path      = "${var.efs_host_mount_path}"
  log_retention_days       = "${var.log_retention_days}"
  monitor_config_file_path = "${var.monitor_config_file_path}"
  monitor_user_agent       = "${var.monitor_user_agent}"
  monitor_timeout          = "${var.monitor_timeout}"
  task_image               = "${var.task_image}"
  task_image_version       = "${var.task_image_version}"
}
