resource "aws_vpc" "vpc" {
  cidr_block = "${var.vpc_cidr}"

  tags {
    Name = "vpc-${var.environment}"
  }
}

resource "aws_internet_gateway" "vpc" {
  vpc_id = "${aws_vpc.vpc.id}"

  tags {
    Environment = "${var.environment}"
  }
}
