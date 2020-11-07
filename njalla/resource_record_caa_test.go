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

func TestAccRecordCAA_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCAACreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCAAExists(
						"njalla_record_caa.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_create",
						"name",
						"testacc1-caa-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_create",
						"content",
						`0 issue "letsencrypt.org"`,
					),
				),
			},
		},
	})
}

func TestAccRecordCAA_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCAAUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCAAExists(
						"njalla_record_caa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update",
						"name",
						"testacc2-caa-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update",
						"content",
						`0 issue "letsencrypt.org"`,
					),
				),
			},
			{
				Config: testAccCheckRecordCAAUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCAAExists(
						"njalla_record_caa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update",
						"name",
						"testacc2-caa-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_update",
						"content",
						`0 iodef "letsencrypt.org"`,
					),
				),
			},
		},
	})
}

func TestAccRecordCAA_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCAAImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCAAExists(
						"njalla_record_caa.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_caa.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordCAA_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordCAAEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordCAAExists(
						"njalla_record_caa.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_caa.test_empty_name",
						"content",
						`0 issue "letsencrypt.org"`,
					),
				),
			},
		},
	})
}

func TestAccRecordCAA_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordCAAInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordCAA_InvalidContent(t *testing.T) {
	expectedErr := regexp.MustCompile("value must follow RFC 8659")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordCAAInvalidContent(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordCAA_InvalidContentFlagTooBig(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"flag must be between 0 and 255: RFC 8659 4.1.1",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordCAADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordCAAInvalidContentFlagTooBig(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordCAADestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_caa" {
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

func testAccCheckRecordCAAExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordCAACreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_create {
  domain = %q
  name = "testacc1-caa-create-name"
  ttl = 10800
  content = "0 issue \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_update {
  domain = %q
  name = "testacc2-caa-update-name1"
  ttl = 10800
  content = "0 issue \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_update {
  domain = %q
  name = "testacc2-caa-update-name2"
  ttl = 3600
  content = "0 iodef \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_import {
  domain = %q
  name = "testacc3-caa-import-name"
  ttl = 10800
  content = "0 issue \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_empty_name {
  domain = %q
  ttl = 10800
  content = "0 issue \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_invalid_ttl {
  domain = %q
  name = "testacc5-caa-invalidttl-name"
  ttl = 999
  content = "0 issue \"letsencrypt.org\""
}
`, domain)
}

func testAccCheckRecordCAAInvalidContent() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_invalid_content {
  domain = %q
  name = "testacc6-caa-invalidcontent-name"
  ttl = 10800
  content = "testacc6-caa-invalidcontent-content"
}
`, domain)
}

func testAccCheckRecordCAAInvalidContentFlagTooBig() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_caa test_flag_too_big {
  domain = %q
  name = "testacc7-caa-flagtoobig-name"
  ttl = 10800
  content = "999 issue \"letsencrypt.org\""
}
`, domain)
}
