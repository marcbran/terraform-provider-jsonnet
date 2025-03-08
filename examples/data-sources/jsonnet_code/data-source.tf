terraform {
  required_providers {
    jsonnet = {
      source = "marcbran/jsonnet"
    }
  }
}

data "jsonnet_code" "example" {
  code = <<EOF
local input = 1;
local output = input + 1;
output
EOF
}

output "result" {
  value = data.jsonnet_code.example.output
}
