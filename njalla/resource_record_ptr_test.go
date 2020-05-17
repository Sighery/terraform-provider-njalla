package njalla

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/Sighery/gonjalla"
)

func TestAccRecordPTR_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordPTRCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordPTRExists(
						"njalla_record_ptr.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_create",
						"name",
						"testacc1-ptr-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_create",
						"content",
						"testacc1-ptr-create-content",
					),
				),
			},
		},
	})
}

func TestAccRecordPTR_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordPTRUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordPTRExists(
						"njalla_record_ptr.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update",
						"name",
						"testacc2-ptr-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update",
						"content",
						"testacc2-ptr-update-content1",
					),
				),
			},
			{
				Config: testAccCheckRecordPTRUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordPTRExists(
						"njalla_record_ptr.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update",
						"name",
						"testacc2-ptr-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_update",
						"content",
						"testacc2-ptr-update-content2",
					),
				),
			},
		},
	})
}

func TestAccRecordPTR_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordPTRImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordPTRExists(
						"njalla_record_ptr.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_ptr.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordPTR_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordPTREmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordPTRExists(
						"njalla_record_ptr.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_ptr.test_empty_name",
						"content",
						"testacc4-ptr-emptyname-content",
					),
				),
			},
		},
	})
}

func TestAccRecordPTR_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordPTRInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordPTRDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_ptr" {
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

func testAccCheckRecordPTRExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordPTRCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_create {
  domain = %q
  name = "testacc1-ptr-create-name"
  ttl = 10800
  content = "testacc1-ptr-create-content"
}
`, domain)
}

func testAccCheckRecordPTRUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_update {
  domain = %q
  name = "testacc2-ptr-update-name1"
  ttl = 10800
  content = "testacc2-ptr-update-content1"
}
`, domain)
}

func testAccCheckRecordPTRUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_update {
  domain = %q
  name = "testacc2-ptr-update-name2"
  ttl = 3600
  content = "testacc2-ptr-update-content2"
}
`, domain)
}

func testAccCheckRecordPTRImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_import {
  domain = %q
  name = "testacc3-ptr-import-name"
  ttl = 10800
  content = "testacc3-ptr-import-content"
}
`, domain)
}

func testAccCheckRecordPTREmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_empty_name {
  domain = %q
  ttl = 10800
  content = "testacc4-ptr-emptyname-content"
}
`, domain)
}

func testAccCheckRecordPTRInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_ptr test_invalid_ttl {
  domain = %q
  name = "testacc5-ptr-invalidttl-name"
  ttl = 999
  content = "testacc5-ptr-invalidttl-content"
}
`, domain)
}
