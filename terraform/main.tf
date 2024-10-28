terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-2"
}

resource "tls_private_key" "generated_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "ec2_key_pair" {
  key_name   = "generated-key"
  public_key = tls_private_key.generated_key.public_key_openssh
}

resource "aws_security_group" "cf_backend_sg" {
  name        = "cf_backend_sg"
  description = "Allow SSH and ports from 1024"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 1024
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "cat_forwarding_backend" {
  ami             = "ami-07c5ecd8498c59db5" # Amazon Linux 2023 AMI
  instance_type   = "t2.micro"
  key_name        = aws_key_pair.ec2_key_pair.key_name
  security_groups = [aws_security_group.cf_backend_sg.name]

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = tls_private_key.generated_key.private_key_pem
    host        = self.public_ip
  }

  provisioner "remote-exec" {
    inline = [
      "echo '${file("~/.ssh/id_rsa.pub")}' >> /home/ec2-user/.ssh/authorized_keys",
      "chmod 600 /home/ec2-user/.ssh/authorized_keys",
    ]
  }

  tags = {
    Name = "cf-backend"
  }
}

resource "aws_eip" "cf_backend_eip" {
  instance = aws_instance.cat_forwarding_backend.id
  vpc      = true
}

output "instance_public_ip" {
  value = aws_eip.cf_backend_eip.public_ip
}
