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

func TestAccRecordA_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordACreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAExists(
						"njalla_record_a.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_create",
						"name",
						"testacc1-a-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_create",
						"content",
						"1.1.1.1",
					),
				),
			},
		},
	})
}

func TestAccRecordA_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAExists(
						"njalla_record_a.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update",
						"name",
						"testacc2-a-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update",
						"content",
						"1.1.1.2",
					),
				),
			},
			{
				Config: testAccCheckRecordAUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAExists(
						"njalla_record_a.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update",
						"name",
						"testacc2-a-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_update",
						"content",
						"1.1.1.3",
					),
				),
			},
		},
	})
}

func TestAccRecordA_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAExists(
						"njalla_record_a.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_a.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordA_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAExists(
						"njalla_record_a.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_a.test_empty_name",
						"content",
						"1.1.1.5",
					),
				),
			},
		},
	})
}

func TestAccRecordA_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordAInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordA_InvalidContent(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected content to contain a valid IPv4 address",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordAInvalidContent(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordADestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_a" {
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

func testAccCheckRecordAExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordACreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_create {
  domain = %q
  name = "testacc1-a-create-name"
  ttl = 10800
  content = "1.1.1.1"
}
`, domain)
}

func testAccCheckRecordAUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_update {
  domain = %q
  name = "testacc2-a-update-name1"
  ttl = 10800
  content = "1.1.1.2"
}
`, domain)
}

func testAccCheckRecordAUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_update {
  domain = %q
  name = "testacc2-a-update-name2"
  ttl = 3600
  content = "1.1.1.3"
}
`, domain)
}

func testAccCheckRecordAImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_import {
  domain = %q
  name = "testacc3-a-import-name"
  ttl = 10800
  content = "1.1.1.4"
}
`, domain)
}

func testAccCheckRecordAEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_empty_name {
  domain = %q
  ttl = 10800
  content = "1.1.1.5"
}
`, domain)
}

func testAccCheckRecordAInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_invalid_ttl {
  domain = %q
  name = "testacc5-a-invalidttl-name"
  ttl = 999
  content = "1.1.1.6"
}
`, domain)
}

func testAccCheckRecordAInvalidContent() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_a test_invalid_content {
  domain = %q
  name = "testacc6-a-invalidcontent-name"
  ttl = 10800
  content = "testacc6-a-invalidcontent-content"
}`, domain)
}
