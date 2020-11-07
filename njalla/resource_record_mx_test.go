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

func TestAccRecordMX_Create(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordMXCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordMXExists(
						"njalla_record_mx.test_create",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_create", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_create",
						"name",
						"testacc1-mx-create-name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_create", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_create", "priority", "10",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_create",
						"content",
						"testacc1-mx-create-content",
					),
				),
			},
		},
	})
}

func TestAccRecordMX_Update(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordMXUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordMXExists(
						"njalla_record_mx.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update",
						"name",
						"testacc2-mx-update-name1",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "priority", "10",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update",
						"content",
						"testacc2-mx-update-content1",
					),
				),
			},
			{
				Config: testAccCheckRecordMXUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordMXExists(
						"njalla_record_mx.test_update",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update",
						"name",
						"testacc2-mx-update-name2",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "ttl", "3600",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update", "priority", "20",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_update",
						"content",
						"testacc2-mx-update-content2",
					),
				),
			},
		},
	})
}

func TestAccRecordMX_Import(t *testing.T) {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordMXImport(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordMXExists(
						"njalla_record_mx.test_import",
					),
				),
			},
			{
				ResourceName:        "njalla_record_mx.test_import",
				ImportStateIdPrefix: fmt.Sprintf("%s:", domain),
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccRecordMX_EmptyName(t *testing.T) {
	// With an empty name field it should get the `DefaultFunc` value `@`
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRecordMXEmptyName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordMXExists(
						"njalla_record_mx.test_empty_name",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_empty_name", "domain", domain,
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_empty_name", "name", "@",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_empty_name", "ttl", "10800",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_empty_name", "priority", "10",
					),
					resource.TestCheckResourceAttr(
						"njalla_record_mx.test_empty_name",
						"content",
						"testacc4-mx-emptyname-content",
					),
				),
			},
		},
	})
}

func TestAccRecordMX_InvalidTTL(t *testing.T) {
	expectedErr := regexp.MustCompile("expected ttl to be one of .+, got 999")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordMXInvalidTTL(),
				ExpectError: expectedErr,
			},
		},
	})
}

func TestAccRecordMX_InvalidPriority(t *testing.T) {
	expectedErr := regexp.MustCompile(
		"expected priority to be one of .+, got 999",
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRecordMXDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckRecordMXInvalidPriority(),
				ExpectError: expectedErr,
			},
		},
	})
}

func testAccCheckRecordMXDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "njalla_record_mx" {
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

func testAccCheckRecordMXExists(resource string) resource.TestCheckFunc {
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

func testAccCheckRecordMXCreate() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_create {
  domain = %q
  name = "testacc1-mx-create-name"
  ttl = 10800
  priority = 10
  content = "testacc1-mx-create-content"
}
`, domain)
}

func testAccCheckRecordMXUpdatePre() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_update {
  domain = %q
  name = "testacc2-mx-update-name1"
  ttl = 10800
  priority = 10
  content = "testacc2-mx-update-content1"
}
`, domain)
}

func testAccCheckRecordMXUpdatePost() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_update {
  domain = %q
  name = "testacc2-mx-update-name2"
  ttl = 3600
  priority = 20
  content = "testacc2-mx-update-content2"
}
`, domain)
}

func testAccCheckRecordMXImport() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_import {
  domain = %q
  name = "testacc3-mx-import-name"
  ttl = 10800
  priority = 10
  content = "testacc3-mx-import-content"
}
`, domain)
}

func testAccCheckRecordMXEmptyName() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_empty_name {
  domain = %q
  ttl = 10800
  priority = 10
  content = "testacc4-mx-emptyname-content"
}
`, domain)
}

func testAccCheckRecordMXInvalidTTL() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_invalid_ttl {
  domain = %q
  name = "testacc5-mx-invalidttl-name"
  ttl = 999
  priority = 10
  content = "testacc5-mx-invalidttl-content"
}
`, domain)
}

func testAccCheckRecordMXInvalidPriority() string {
	domain := os.Getenv("NJALLA_TESTACC_DOMAIN")
	return fmt.Sprintf(`
resource njalla_record_mx test_invalid_priority {
  domain = %q
  name = "testacc6-mx-invalidpriority-name"
  ttl = 10800
  priority = 999
  content = "testacc6-mx-invalidpriority-content"
}
`, domain)
}
