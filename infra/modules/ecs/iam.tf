resource "aws_iam_instance_profile" "cluster_instance" {
  name = "${var.cluster_name}-${var.environment}"
  role = "${aws_iam_role.cluster_instance.id}"
}

resource "aws_iam_role" "cluster_instance" {
  name = "${var.cluster_name}-${var.environment}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["ec2.amazonaws.com"]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

# Ref https://docs.aws.amazon.com/AmazonECS/latest/developerguide/instance_IAM_role.html.
# For EFS - https://docs.aws.amazon.com/efs/latest/ug/access-control-managing-permissions.html
resource "aws_iam_role_policy" "cluster_instance" {
  name = "${var.cluster_name}-${var.environment}"
  role = "${aws_iam_role.cluster_instance.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:CreateCluster",
        "ecs:DeregisterContainerInstance",
        "ecs:DiscoverPollEndpoint",
        "ecs:Poll",
        "ecs:RegisterContainerInstance",
        "ecs:StartTelemetrySession",
        "ecs:Submit*",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    },
    {
      "Sid" : "AllowFileSystemPermissions",  
      "Effect": "Allow",
      "Action": [
        "elasticfilesystem:CreateFileSystem",
        "elasticfilesystem:CreateMountTarget"
      ],
      "Resource": "arn:aws:elasticfilesystem:us-west-2:account-id:file-system/*"
    },
    {
     "Sid" : "AllowEC2Permissions",
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeSubnets",
        "ec2:CreateNetworkInterface",
        "ec2:DescribeNetworkInterfaces"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}
