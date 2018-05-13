resource "aws_cloudwatch_log_group" "log_group" {
  name              = "${var.name}-${var.environment}"
  retention_in_days = "${var.log_retention_days}"
}
