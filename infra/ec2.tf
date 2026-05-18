
resource "aws_instance" "public_instance" {
  ami                    = "ami-091138d0f0d41ff90"
  instance_type          = "t2.micro"
  subnet_id              = aws_subnet.public_subnet.id
  key_name               = data.aws_key_pair.key.key_name
  vpc_security_group_ids = [aws_security_group.proaula_security_G.id]

  tags = {
    "Name" = "Proaula"
  }
}
