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

func TestAccRecordNS_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNSCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNSExists(
						"njalla_record_ns.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_create",
						"name",
						"testacc1-ns-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_create",
						"content",
						"testacc1-ns-create-content",
					),
				),
			},
		},
	})
}

func TestAccRecordNS_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNSUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNSExists(
						"njalla_record_ns.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update",
						"name",
						"testacc2-ns-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update",
						"content",
						"testacc2-ns-update-content1",
					),
				),
			},
			{
				Config: testAccCheckRecordNSUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNSExists(
						"njalla_record_ns.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update",
						"name",
						"testacc2-ns-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ns.test_update",
						"content",
						"testacc2-ns-update-content2",
					),
				),
			},
		},
	})
}

func TestAccRecordNS_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNSImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNSExists(
						"njalla_record_ns.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_ns.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordNS_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNSDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordNSInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordNSDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_ns" {
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
			strID := fmt.Sprint(record.ID)
			if strID == rs.Primary.ID {
				return fmt.Errorf(
					"Record %s still exists in domain %s",
					rs.Primary.ID, domain,
				)
			}
		}
	}

	return nil
}

func testAccCheckRecordNSExists(resource string) resource.TestCheckFunc {
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
			strID := fmt.Sprint(record.ID)
			if strID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf(
			"Record %s doesn't exist for domain %s", rs.Primary.ID, domain,
		)
	}
}

func testAccCheckRecordNSCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ns test_create {
  domain = %q
  name = "testacc1-ns-create-name"
  ttl = 10800
  content = "testacc1-ns-create-content"
}
`, domain)
}

func testAccCheckRecordNSUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ns test_update {
  domain = %q
  name = "testacc2-ns-update-name1"
  ttl = 10800
  content = "testacc2-ns-update-content1"
}
`, domain)
}

func testAccCheckRecordNSUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ns test_update {
  domain = %q
  name = "testacc2-ns-update-name2"
  ttl = 3600
  content = "testacc2-ns-update-content2"
}
`, domain)
}

func testAccCheckRecordNSImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ns test_import {
  domain = %q
  name = "testacc3-ns-import-name"
  ttl = 10800
  content = "testacc3-ns-import-content"
}
`, domain)
}

func testAccCheckRecordNSInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ns test_invalid_ttl {
  domain = %q
  name = "testacc4-ns-invalidttl-name"
  ttl = 999
  content = "testacc4-ns-invalidttl-content"
}
`, domain)
}
