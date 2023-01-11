resource "null_resource" "example1" {
  triggers = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    command = " while true; do echo 'Hit CTRL+C'; sleep 1; done"
  }
}
