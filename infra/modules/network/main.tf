module "vpc" {
  source      = "../vpc"
  vpc_cidr    = "${var.vpc_cidr}"
  environment = "${var.environment}"
}

module "private_subnet" {
  source             = "../subnet"
  name               = "private-subnet"
  environment        = "${var.environment}"
  vpc_id             = "${module.vpc.id}"
  subnet_cidrs       = "${var.private_subnet_cidrs}"
  availability_zones = "${var.availability_zones}"
}

module "public_subnet" {
  source             = "../subnet"
  name               = "public-subnet"
  environment        = "${var.environment}"
  vpc_id             = "${module.vpc.id}"
  subnet_cidrs       = "${var.public_subnet_cidrs}"
  availability_zones = "${var.availability_zones}"
}

module "nat_gateway" {
  source     = "../nat_gateway"
  subnet_ids = "${module.public_subnet.subnet_ids}"
  count      = "${length(var.public_subnet_cidrs)}"
}

resource "aws_route" "public_igw_route" {
  count                  = "${length(var.public_subnet_cidrs)}"
  route_table_id         = "${element(module.public_subnet.route_table_ids, count.index)}"
  gateway_id             = "${module.vpc.igw}"
  destination_cidr_block = "${var.destination_cidr_block}"
}

resource "aws_route" "nat_gateway_route" {
  route_table_id         = "${element(module.private_subnet.route_table_ids, count.index)}"
  nat_gateway_id         = "${element(module.nat_gateway.ids, count.index)}"
  destination_cidr_block = "${var.destination_cidr_block}"
  count                  = "${length(var.private_subnet_cidrs)}"
}

// NAT Gateway takes some time to come up. The below statement injects a dependency and we wait till NAT Gatewcay is ready
resource "null_resource" "dummy_dependency" {
  depends_on = ["module.nat_gateway"]
}
