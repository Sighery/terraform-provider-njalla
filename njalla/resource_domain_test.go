package njalla

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/Sighery/gonjalla"
)

func TestAccDomain_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	func() { testAccPreCheck(t) },
		Providers:	testAccProviders,
		CheckDestroy:	testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDomainCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(
						"njalla_domain.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_domain.test_create",
						"name",
						domain,
					),
				),
			},
		},
	})
}

func TestAccDomain_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	func() { testAccPreCheck(t) },
		Providers:	testAccProviders,
		CheckDestroy:	testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDomainImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists(
						"njalla_domain.test_import",
					),
				),
			},
			{
				ResourceName:	"njalla_server.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:	true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDomain_UnavaliableDomain(t *testing.T) {
	expectedErr := regexp.MustCompile(`Failed to buy \b(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}\b`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	func() { testAccPreCheck(t) },
		Providers:	testAccProviders,
		CheckDestroy:	testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDomainUnavaliableDomain(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckDomainCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_domain test_create {
  name = %q
  years = 1
}
`, domain)
}

func testAccCheckDomainImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_domain test_import {
  name = %q
  years = 1
}
`, domain)
}

func testAccCheckDomainUnavaliableDomain() string {
	return fmt.Sprintf(`
resource njalla_domain test_import {
  name = "example.com"
  years = 1
}`)
}

func testAccCheckDomainExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		domainName := os.Getenv("NJALLA_TESTACC_DOMAIN")
		domains, err := gonjalla.ListDomains(config.Token)
		if err != nil {
			return fmt.Errorf(
				"Error fetching domain data %s: %s",
				domainName, err,
			)
		}

		for _, domain := range domains {
			if domain.Name == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf(
			"%s doesn't exist", domainName,
		)
	}
}

func testAccCheckDomainDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domainName := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_domain" {
			continue
		}

		domains, err := gonjalla.ListDomains(config.Token)
		if err != nil {
			return fmt.Errorf(
				"Error fetching domain data %s: %s",
				domainName, err,
			)
		}

		for _, domain := range domains {
			if domain.Name == rs.Primary.ID {
				return fmt.Errorf(
					"%s still exists",
					domainName,
				)
			}
		}
	}

	return nil
}