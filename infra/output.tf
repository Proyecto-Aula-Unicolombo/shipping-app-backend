output "ec2_public_ip" {
  description = "public ip of instance"
  value       = aws_instance.public_instance.public_ip
}
