resource "aws_ecs_service" "service" {
  name            = "${var.name}-${var.environment}"
  cluster         = "${var.cluster_id}"
  task_definition = "${aws_ecs_task_definition.task.arn}"
  desired_count   = 1
}

resource "aws_ecs_task_definition" "task" {
  family = "${var.name}-${var.environment}"

  volume {
    name      = "monitor"
    host_path = "${var.efs_host_mount_path}"
  }

  container_definitions = <<EOF
[
    {
        "name": "${var.name}-${var.environment}",
        "image": "${var.task_image}:${var.task_image_version}",
        "cpu": 10,
        "memoryReservation": 250,
        "mountPoints": [
            {
                "sourceVolume": "monitor",
                "containerPath": "${var.efs_task_mount_path}",
                "readOnly": false
            }
        ],

        "logConfiguration": {
            "logDriver": "awslogs",
             "options": {
                    "awslogs-group": "${aws_cloudwatch_log_group.log_group.name}",
                    "awslogs-region": "${var.region}",
                    "awslogs-stream-prefix": "${var.name}-${var.environment}"
                }
        },
        "environment" :[
            {
                "name": "MONITOR_CONFIG_FILE_PATH",
                "value": "${var.monitor_config_file_path}"
            },
            {
                "name": "MONITOR_LOG_DIR_PATH",
                "value": "${var.efs_task_mount_path}"
            },
            {
                "name": "MONITOR_USER_AGENT",
                "value": "${var.monitor_user_agent}"
            },
            {
                "name": "MONITOR_TIME_OUT",
                "value": "${var.monitor_timeout}"
            }
        ]
    }
]
EOF
}
