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

func TestAccRecordTLSA_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTLSACreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTLSAExists(
						"njalla_record_tlsa.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_create",
						"name",
						"testacc1-tlsa-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_create",
						"content",
						"3 0 0 30820307308201efa003020102020",
					),
				),
			},
		},
	})
}

func TestAccRecordTLSA_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTLSAUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTLSAExists(
						"njalla_record_tlsa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update",
						"name",
						"testacc2-tlsa-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update",
						"content",
						"0 0 1 d2abde240d7cd3ee6b4b28c54df034b9",
					),
				),
			},
			{
				Config: testAccCheckRecordTLSAUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTLSAExists(
						"njalla_record_tlsa.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update",
						"name",
						"testacc2-tlsa-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_update",
						"content",
						"0 0 1 7983a1d16e8a410e4561cb106618e971",
					),
				),
			},
		},
	})
}

func TestAccRecordTLSA_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTLSAImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTLSAExists(
						"njalla_record_tlsa.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_tlsa.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordTLSA_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordTLSAEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordTLSAExists(
						"njalla_record_tlsa.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_tlsa.test_empty_name",
						"content",
						"1 1 2 92003ba34942dc74152e2f2c408d29ec",
					),
				),
			},
		},
	})
}

func TestAccRecordTLSA_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_InvalidContent(t *testing.T) {
	expectedErr := regexp.MustCompile(
		`expected 4 arguments, .* Check RFC 6698`,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAInvalidContent(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_StringCertificateUsage(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Certificate Usage field to be int",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAStringCertificateUsage(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_InvalidIntCertificateUsage(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Certificate Usage field to be between 0 and 255",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAInvalidIntCertificateUsage(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_StringSelector(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Selector field to be int",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAStringSelector(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_InvalidIntSelector(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Selector field to be between 0 and 255",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAInvalidIntSelector(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_StringMatchingType(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Matching Type field to be int",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAStringMatchingType(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordTLSA_InvalidIntMatchingType(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected Matching Type field to be between 0 and 255",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordTLSADestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordTLSAInvalidIntMatchingType(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordTLSADestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_tlsa" {
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

func testAccCheckRecordTLSAExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordTLSACreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_create {
  domain = %q
  name = "testacc1-tlsa-create-name"
  ttl = 10800
  content = "3 0 0 30820307308201efa003020102020"
}
`, domain)
}

func testAccCheckRecordTLSAUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_update {
  domain = %q
  name = "testacc2-tlsa-update-name1"
  ttl = 10800
  content = "0 0 1 d2abde240d7cd3ee6b4b28c54df034b9"
}
`, domain)
}

func testAccCheckRecordTLSAUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_update {
  domain = %q
  name = "testacc2-tlsa-update-name2"
  ttl = 3600
  content = "0 0 1 7983a1d16e8a410e4561cb106618e971"
}
`, domain)
}

func testAccCheckRecordTLSAImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_import {
  domain = %q
  name = "testacc3-tlsa-import-name"
  ttl = 10800
  content = "1 1 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_empty_name {
  domain = %q
  ttl = 10800
  content = "1 1 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_invalid_ttl {
  domain = %q
  name = "testacc5-tlsa-invalidttl-name"
  ttl = 999
  content = "1 1 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAInvalidContent() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_invalid_content {
  domain = %q
  name = "testacc6-tlsa-invalidcontent-name"
  ttl = 10800
  content = "testacc6-tlsa-invalidcontent-content"
}
`, domain)
}

func testAccCheckRecordTLSAStringCertificateUsage() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_string_certificate_usage {
  domain = %q
  name = "testacc7-tlsa-stringcertificateusage-name"
  ttl = 10800
  content = "test 1 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAInvalidIntCertificateUsage() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_invalid_int_certificate_usage {
  domain = %q
  name = "testacc8-tlsa-invalidintcertificateusage-name"
  ttl = 10800
  content = "999 1 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAStringSelector() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_string_selector {
  domain = %q
  name = "testacc9-tlsa-stringselector-name"
  ttl = 10800
  content = "0 test 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAInvalidIntSelector() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_invalid_int_selector {
  domain = %q
  name = "testacc10-tlsa-invalidintselector-name"
  ttl = 10800
  content = "0 999 2 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAStringMatchingType() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_string_matching_type {
  domain = %q
  name = "testacc11-tlsa-stringmatchingtype-name"
  ttl = 10800
  content = "0 1 test 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}

func testAccCheckRecordTLSAInvalidIntMatchingType() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_tlsa test_invalid_int_matching_type {
  domain = %q
  name = "testacc12-tlsa-invalidintmatchingtype-name"
  ttl = 10800
  content = "0 1 999 92003ba34942dc74152e2f2c408d29ec"
}
`, domain)
}
