resource "aws_ecs_cluster" "cluster" {
  name = "${var.cluster_name}-${var.environment}"
}

resource "aws_autoscaling_group" "cluster" {
  name                 = "${var.cluster_name}-${var.environment}"
  min_size             = "${var.min_size}"
  max_size             = "${var.max_size}"
  desired_capacity     = "${var.desired_capacity}"
  vpc_zone_identifier  = ["${var.private_subnet_ids}"]
  launch_configuration = "${aws_launch_configuration.cluster.id}"

  tag {
    key                 = "Network_Depends_ID"
    value               = "${var.depends_id}"
    propagate_at_launch = "false"
  }
}

resource "aws_launch_configuration" "cluster" {
  name_prefix          = "${var.cluster_name}-${var.environment}-"
  instance_type        = "${var.instance_type}"
  image_id             = "${var.aws_ami}"
  security_groups      = ["${aws_security_group.cluster.id}"]
  iam_instance_profile = "${aws_iam_instance_profile.cluster_instance.id}"

  user_data = <<EOF
#!/bin/bash
sudo yum update -y
sudo yum install -y nfs-utils


echo "ECS_CLUSTER=${aws_ecs_cluster.cluster.name}"  >> /etc/ecs/ecs.config
echo 'ECS_AVAILABLE_LOGGING_DRIVERS=["awslogs"]'  >> /etc/ecs/ecs.config


sudo mkdir ${var.efs_mount_path}
sudo mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 $(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone).${var.efs_id}.efs.${var.region}.amazonaws.com:/ ${var.efs_mount_path}
sudo chmod 777 ${var.efs_mount_path}

sudo stop ecs


# docker stop ecs-agent
# docker rm ecs-agent


# docker run --name ecs-agent \
# --detach=true \
# --restart=on-failure:10 \
# --volume=/var/run:/var/run \
# --volume=/var/log/ecs/:/log \
# --volume=/var/lib/ecs/data:/data \
# --volume=/etc/ecs:/etc/ecs \
# --net=host \
# --env=ECS_DATADIR=/data \
# --env=ECS_ENABLE_TASK_IAM_ROLE=true \
# --env=ECS_ENABLE_TASK_IAM_ROLE_NETWORK_HOST=true \
# --env=ECS_LOGFILE=/log/ecs-agent.log \
# --env=ECS_AVAILABLE_LOGGING_DRIVERS=["awslogs"] \
# --env=ECS_LOGLEVEL=info \
# --env=ECS_CLUSTER=${aws_ecs_cluster.cluster.name} \
# --env=ECS_APPARMOR_CAPABLE=true \
# --log-driver=awslogs \
# --log-opt=awslogs-region=${var.region} \
# --log-opt=awslogs-group=${var.cluster_name} \
# amazon/amazon-ecs-agent:latest

sudo service docker restart
sudo start ecs
EOF

  lifecycle {
    create_before_destroy = true
  }
}
