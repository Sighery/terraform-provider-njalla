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

func TestAccRecordAAAA_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAAAACreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAAAAExists(
						"njalla_record_aaaa.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_create",
						"name",
						"testacc1-aaaa-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_create",
						"content",
						"2001:db8::8a2e:370:7331",
					),
				),
			},
		},
	})
}

func TestAccRecordAAAA_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAAAAUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAAAAExists(
						"njalla_record_aaaa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update",
						"name",
						"testacc2-aaaa-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update",
						"content",
						"2001:db8::8a2e:370:7332",
					),
				),
			},
			{
				Config: testAccCheckRecordAAAAUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAAAAExists(
						"njalla_record_aaaa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update",
						"name",
						"testacc2-aaaa-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_update",
						"content",
						"2001:db8::8a2e:370:7333",
					),
				),
			},
		},
	})
}

func TestAccRecordAAAA_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAAAAImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAAAAExists(
						"njalla_record_aaaa.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_aaaa.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordAAAA_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordAAAAEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordAAAAExists(
						"njalla_record_aaaa.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_aaaa.test_empty_name",
						"content",
						"2001:db8::8a2e:370:7335",
					),
				),
			},
		},
	})
}

func TestAccRecordAAAA_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordAAAAInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordAAAA_InvalidContent(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected content to contain a valid IPv6 address",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordAAAADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordAAAAInvalidContent(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordAAAADestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_aaaa" {
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

func testAccCheckRecordAAAAExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordAAAACreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_create {
  domain = %q
  name = "testacc1-aaaa-create-name"
  ttl = 10800
  content = "2001:db8::8a2e:370:7331"
}
`, domain)
}

func testAccCheckRecordAAAAUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_update {
  domain = %q
  name = "testacc2-aaaa-update-name1"
  ttl = 10800
  content = "2001:db8::8a2e:370:7332"
}
`, domain)
}

func testAccCheckRecordAAAAUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_update {
  domain = %q
  name = "testacc2-aaaa-update-name2"
  ttl = 3600
  content = "2001:db8::8a2e:370:7333"
}
`, domain)
}

func testAccCheckRecordAAAAImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_import {
  domain = %q
  name = "testacc3-aaaa-import-name"
  ttl = 10800
  content = "2001:db8::8a2e:370:7334"
}
`, domain)
}

func testAccCheckRecordAAAAEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_empty_name {
  domain = %q
  ttl = 10800
  content = "2001:db8::8a2e:370:7335"
}
`, domain)
}

func testAccCheckRecordAAAAInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_invalid_ttl {
  domain = %q
  name = "testacc5-aaaa-invalidttl-name"
  ttl = 999
  content = "2001:db8::8a2e:370:7336"
}
`, domain)
}

func testAccCheckRecordAAAAInvalidContent() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_aaaa test_invalid_content {
  domain = %q
  name = "testacc6-aaaa-invalidcontent-name"
  ttl = 10800
  content = "testacc6-aaaa-invalidcontent-content"
}`, domain)
}
