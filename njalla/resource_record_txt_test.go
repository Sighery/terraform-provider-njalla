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

func TestAccRecordTXT_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTXTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTXTCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTXTExists(
						"njalla_record_txt.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_create",
						"name",
						"testacc1-txt-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_create",
						"content",
						"testacc1-txt-create-content",
					),
				),
			},
		},
	})
}

func TestAccRecordTXT_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTXTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTXTUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTXTExists(
						"njalla_record_txt.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update",
						"name",
						"testacc2-txt-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update",
						"content",
						"testacc2-txt-update-content1",
					),
				),
			},
			{
				Config: testAccCheckRecordTXTUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTXTExists(
						"njalla_record_txt.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update",
						"name",
						"testacc2-txt-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_update",
						"content",
						"testacc2-txt-update-content2",
					),
				),
			},
		},
	})
}

func TestAccRecordTXT_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTXTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTXTImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTXTExists(
						"njalla_record_txt.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_txt.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordTXT_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTXTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTXTEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTXTExists(
						"njalla_record_txt.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_txt.test_empty_name",
						"content",
						"testacc4-txt-emptyname-content",
					),
				),
			},
		},
	})
}

func TestAccRecordTXT_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTXTDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTXTInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordTXTDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_txt" {
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

func testAccCheckRecordTXTExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordTXTCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_create {
  domain = %q
  name = "testacc1-txt-create-name"
  ttl = 10800
  content = "testacc1-txt-create-content"
}
`, domain)
}

func testAccCheckRecordTXTUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_update {
  domain = %q
  name = "testacc2-txt-update-name1"
  ttl = 10800
  content = "testacc2-txt-update-content1"
}
`, domain)
}

func testAccCheckRecordTXTUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_update {
  domain = %q
  name = "testacc2-txt-update-name2"
  ttl = 3600
  content = "testacc2-txt-update-content2"
}
`, domain)
}

func testAccCheckRecordTXTImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_import {
  domain = %q
  name = "testacc3-txt-import-name"
  ttl = 10800
  content = "testacc3-txt-import-content"
}
`, domain)
}

func testAccCheckRecordTXTEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_empty_name {
  domain = %q
  ttl = 10800
  content = "testacc4-txt-emptyname-content"
}
`, domain)
}

func testAccCheckRecordTXTInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_txt test_invalid_ttl {
  domain = %q
  name = "testacc5-txt-invalidttl-name"
  ttl = 999
  content = "testacc5-txt-invalidttl-content"
}
`, domain)
}
