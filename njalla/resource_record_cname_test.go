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

func TestAccRecordCNAME_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCNAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCNAMECreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCNAMEExists(
						"njalla_record_cname.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_create",
						"name",
						"testacc1-cname-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_create",
						"content",
						"testacc1-cname-create-content.com",
					),
				),
			},
		},
	})
}

func TestAccRecordCNAME_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCNAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCNAMEUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCNAMEExists(
						"njalla_record_cname.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update",
						"name",
						"testacc2-cname-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update",
						"content",
						"testacc2-cname-update-content1.com",
					),
				),
			},
			{
				Config: testAccCheckRecordCNAMEUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCNAMEExists(
						"njalla_record_cname.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update",
						"name",
						"testacc2-cname-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_update",
						"content",
						"testacc2-cname-update-content2.com",
					),
				),
			},
		},
	})
}

func TestAccRecordCNAME_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCNAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCNAMEImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCNAMEExists(
						"njalla_record_cname.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_cname.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordCNAME_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCNAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCNAMEEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCNAMEExists(
						"njalla_record_cname.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_cname.test_empty_name",
						"content",
						"testacc4-cname-emptyname-content.com",
					),
				),
			},
		},
	})
}

func TestAccRecordCNAME_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCNAMEDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordCNAMEInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordCNAMEDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_cname" {
			continue
		}

		records, err := gonjalla.ListRecords(config.Token, domain)
		if err != nil {
			return fmt.Errorf(
				"Error fetching the records data for domain %s: %s",
				domain, err,
			)
		}

		for _, record := range records {
			if record.ID == rs.Primary.ID {
				return fmt.Errorf(
					"Record %s still exists in domain %s",
					rs.Primary.ID, domain,
				)
			}
		}
	}

	return nil
}

func testAccCheckRecordCNAMEExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
		records, err := gonjalla.ListRecords(config.Token, domain)
		if err != nil {
			return fmt.Errorf(
				"Error fetching the records data for domain %s: %s",
				domain, err,
			)
		}

		for _, record := range records {
			if record.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf(
			"Record %s doesn't exist for domain %s", rs.Primary.ID, domain,
		)
	}
}

func testAccCheckRecordCNAMECreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_create {
  domain = %q
  name = "testacc1-cname-create-name"
  ttl = 10800
  content = "testacc1-cname-create-content.com"
}
`, domain)
}

func testAccCheckRecordCNAMEUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_update {
  domain = %q
  name = "testacc2-cname-update-name1"
  ttl = 10800
  content = "testacc2-cname-update-content1.com"
}
`, domain)
}

func testAccCheckRecordCNAMEUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_update {
  domain = %q
  name = "testacc2-cname-update-name2"
  ttl = 3600
  content = "testacc2-cname-update-content2.com"
}
`, domain)
}

func testAccCheckRecordCNAMEImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_import {
  domain = %q
  name = "testacc3-cname-import-name"
  ttl = 10800
  content = "testacc3-cname-import-content.com"
}
`, domain)
}

func testAccCheckRecordCNAMEEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_empty_name {
  domain = %q
  ttl = 10800
  content = "testacc4-cname-emptyname-content.com"
}
`, domain)
}

func testAccCheckRecordCNAMEInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_cname test_invalid_ttl {
  domain = %q
  name = "testacc5-cname-invalidttl-name"
  ttl = 999
  content = "testacc5-cname-invalidttl-content.com"
}
`, domain)
}
