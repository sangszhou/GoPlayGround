variable "list_example" {
  description = "An example of a list in Terraform" type = "list"
  default = [1, 2, 3]
}

variable "map_example" {
  description = "An example of a map in Terraform" type = "map"
  default = {
    key1 = "value1"
    key2 = "value2"
    key3 = "value3"
  }
}

variable "server_port" {
  description = "The port the server will use for HTTP requests"
  default = 8080
}


provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "example" {
//  ami The Amazon Machine Image (AMI) to run on the EC2 Instance
  ami = "ami-40d28157"
//  t2.micro, which has 1 virtual CPU, 1GB of memory
  instance_type = "t2.micro"

  //指定的数据，启动 vm 后执行
  user_data = <<-EOF
    #!/bin/bash
    echo "Hello, World" > index.html
    nohup busybox httpd -f -p 8080 &
    EOF

  tags {
    Name = "terraform-example"
  }
}

resource "aws_security_group" "instance" {
  name = "terraform-example-instance"
  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}