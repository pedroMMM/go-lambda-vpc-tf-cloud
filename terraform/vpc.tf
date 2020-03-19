resource "aws_vpc" "test_vpc" {
  count = var.vpc_count

  # I am feeling very dirty by having this crazy CIDR but this VPC
  # will never house anything.
  cidr_block = "10.0.0.0/16"

  tags = {
    Name  = "Test VPC #${count.index}"
    App   = local.app
    Owner = var.owner
  }
}
