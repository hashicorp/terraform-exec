resource "aws_instance" "web" {
  ami           = "${data.aws_ami.amazon.id}"
  instance_type = "t3.micro"
  count         = 2

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_elb" "web" {
  instances = ["${aws_instance.web.*.id}"]
  subnets   = ["${aws_subnet.test.*.id}"]
  listener {
    instance_port     = 8000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}

data "aws_ami" "amazon" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-*-x86_64-gp2"]
  }
}
