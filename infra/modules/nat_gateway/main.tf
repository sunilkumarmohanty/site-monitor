resource "aws_nat_gateway" "gateway" {
  allocation_id = "${element(aws_eip.gateway.*.id, count.index)}"
  subnet_id     = "${element(var.subnet_ids, count.index)}"
  count         = "${var.count}"
}

resource "aws_eip" "gateway" {
  vpc   = true
  count = "${var.count}"
}
