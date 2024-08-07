terraform {
  required_providers {
    jsonnet = {
      source = "marcbran/jsonnet"
    }
  }
}

output "output" {
  value = jsondecode(provider::jsonnet::evaluate(<<-EOF
{
  person1: {
    name: "Alice",
    welcome: "Hello " + self.name + "!",
  },
  person2: self.person1 { name: "Bob" },
}
EOF
  ))["person2"]["name"]
}
