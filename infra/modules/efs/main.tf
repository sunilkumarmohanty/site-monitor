resource "aws_efs_file_system" "efs" {
  tags {
    Name = "${var.name}-${var.environment}"
  }

  encrypted = "true"
}

resource "aws_efs_mount_target" "mount" {
  count           = "${var.count}"
  subnet_id       = "${element(var.subnet_ids, count.index)}"
  security_groups = ["${aws_security_group.efsmount.id}"]
  file_system_id  = "${aws_efs_file_system.efs.id}"
}
