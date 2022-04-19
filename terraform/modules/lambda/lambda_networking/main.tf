resource "aws_nat_gateway" "nat_gateway" {
  subnet_id     = var.public_subnet_id
  allocation_id = aws_eip.nat.allocation_id
}


resource "aws_route_table_association" "public" {
  subnet_id      = var.public_subnet_id
  route_table_id = var.public_route_table_id
}

resource "aws_route_table" "private" {
  vpc_id = var.vpc_id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway.id
  }
}

resource "aws_route_table_association" "private" {
  subnet_id      = var.private_subnet_id
  route_table_id = aws_route_table.private.id
}

resource "aws_eip" "nat" {
  vpc = true
}
