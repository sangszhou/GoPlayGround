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

resource "aws_security_group" "instance" {
  name = "terraform-example-instance"
  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_launch_configuration" "example" {
  image_id = "ami-40d28157"
  instance_type = "t2.micro"
  security_groups = ["${aws_security_group.instance.id}"]
  user_data = <<-EOF
  #!/bin/bash
    echo "Hello, World" > index.html
    nohup busybox httpd -f -p "${var.server_port}" &
    EOF
  lifecycle {
    create_before_destroy = true
  }
}