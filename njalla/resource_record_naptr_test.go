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

func TestAccRecordNAPTR_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNAPTRCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNAPTRExists(
						"njalla_record_naptr.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_create",
						"name",
						"testacc1-naptr-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_create",
						"content",
						`100 10 "" "" "/urn:cid:.+@([^\.]+\.)(.*)$/\2/i" .`,
					),
				),
			},
		},
	})
}

func TestAccRecordNAPTR_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNAPTRUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNAPTRExists(
						"njalla_record_naptr.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update",
						"name",
						"testacc2-naptr-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update",
						"content",
						`100 10 "" "" "/urn:cid:.+@([^\.]+\.)(.*)$/\2/i" .`,
					),
				),
			},
			{
				Config: testAccCheckRecordNAPTRUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNAPTRExists(
						"njalla_record_naptr.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update",
						"name",
						"testacc2-naptr-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_update",
						"content",
						`100 50 "S" "z3950+I2L+I2C" "" _z3950._tcp.gatech.edu.`,
					),
				),
			},
		},
	})
}

func TestAccRecordNAPTR_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNAPTRImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNAPTRExists(
						"njalla_record_naptr.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_naptr.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordNAPTR_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordNAPTREmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordNAPTRExists(
						"njalla_record_naptr.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_naptr.test_empty_name",
						"content",
						`100 50 "S" "http+I2L+I2C+I2R" "" _http._tcp.gatech.edu.`,
					),
				),
			},
		},
	})
}

func TestAccRecordNAPTR_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordNAPTRInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordNAPTR_InvalidContent(t *testing.T) {
	expectedErr := regexp.MustCompile(
		`expected 6\+ arguments, .* Check RFC 2915`,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordNAPTRInvalidContent(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordNAPTR_StringOrder(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Order field to be int",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordNAPTRStringOrder(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordNAPTR_StringPreference(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Preference field to be int",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordNAPTRDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordNAPTRStringPreference(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordNAPTRDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_naptr" {
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

func testAccCheckRecordNAPTRExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordNAPTRCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_create {
  domain = %q
  name = "testacc1-naptr-create-name"
  ttl = 10800
  content = "100 10 \"\" \"\" \"/urn:cid:.+@([^\\.]+\\.)(.*)$/\\2/i\" ."
}
`, domain)
}

func testAccCheckRecordNAPTRUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_update {
  domain = %q
  name = "testacc2-naptr-update-name1"
  ttl = 10800
  content = "100 10 \"\" \"\" \"/urn:cid:.+@([^\\.]+\\.)(.*)$/\\2/i\" ."
}
`, domain)
}

func testAccCheckRecordNAPTRUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_update {
  domain = %q
  name = "testacc2-naptr-update-name2"
  ttl = 3600
  content = "100 50 \"S\" \"z3950+I2L+I2C\" \"\" _z3950._tcp.gatech.edu."
}
`, domain)
}

func testAccCheckRecordNAPTRImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_import {
  domain = %q
  name = "testacc3-naptr-import-name"
  ttl = 10800
  content = "100 50 \"S\" \"rcds+I2C\" \"\" _rcds._udp.gatech.edu."
}
`, domain)
}

func testAccCheckRecordNAPTREmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_empty_name {
  domain = %q
  ttl = 10800
  content = "100 50 \"S\" \"http+I2L+I2C+I2R\" \"\" _http._tcp.gatech.edu."
}
`, domain)
}

func testAccCheckRecordNAPTRInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_invalid_ttl {
  domain = %q
  name = "testacc5-naptr-invalidttl-name"
  ttl = 999
  content = "100 50 \"S\" \"http+I2L+I2C+I2R\" \"\" _http._tcp.gatech.edu."
}
`, domain)
}

func testAccCheckRecordNAPTRInvalidContent() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_invalid_content {
  domain = %q
  name = "testacc6-naptr-invalidcontent-name"
  ttl = 10800
  content = "testacc6-naptr-invalidcontent-content"
}
`, domain)
}

func testAccCheckRecordNAPTRStringOrder() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_string_order {
  domain = %q
  name = "testacc7-naptr-stringorder-name"
  ttl = 10800
  content = "test 50 \"S\" \"http+I2L+I2C+I2R\" \"\" _http._tcp.gatech.edu."
}
`, domain)
}

func testAccCheckRecordNAPTRStringPreference() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_naptr test_string_preference {
  domain = %q
  name = "testacc8-naptr-stringpreference-name"
  ttl = 10800
  content = "100 test \"S\" \"http+I2L+I2C+I2R\" \"\" _http._tcp.gatech.edu."
}
`, domain)
}
