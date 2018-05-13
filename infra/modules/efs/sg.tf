resource "aws_security_group" "efsmount" {
  name   = "${var.name}-${var.environment}"
  vpc_id = "${var.vpc_id}"

  ingress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
