resource "aws_vpc" "virginia_test_vpc" {
  cidr_block = var.virginia_vpc.cidr

  tags = {
    "Name" = "vpc_virginia"
  }
}


resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.virginia_test_vpc.id
  cidr_block              = var.subnets[0]
  map_public_ip_on_launch = true

  tags = {
    "Name" = "public_subnet"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.virginia_test_vpc.id

  tags = {
    "Name" = "proaula_igw"
  }
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.virginia_test_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    "Name" = "public_route_table"
  }
}

resource "aws_route_table_association" "public_rt_assoc" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_subnet" "private_subnet" {
  vpc_id     = aws_vpc.virginia_test_vpc.id
  cidr_block = var.subnets[1]

  tags = {
    "Name" = "private_subnet"
  }

  depends_on = [
    aws_subnet.public_subnet // depends to create first the public subnet

  ]
}
