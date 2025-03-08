package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCodeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCodeDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.jsonnet_code.test",
						tfjsonpath.New("output"),
						knownvalue.StringExact("2"),
					),
				},
			},
		},
	})
}

func TestAccInvalidCodeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInvalidCodeDataSourceConfig(),
				ExpectError: regexp.MustCompile("Unable to evaluate Jsonnet"),
			},
		},
	})
}

func testAccCodeDataSourceConfig() string {
	return `
data "jsonnet_code" "test" {
  code = <<EOF
local input = 1;
local output = input + 1;
output
EOF
}
`
}

func testAccInvalidCodeDataSourceConfig() string {
	return `
data "jsonnet_code" "test" {
  code = <<EOF
local output = {
  foo: "bar"
}
EOF
}
`
}
