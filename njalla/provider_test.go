package njalla

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"njalla": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NJALLA_API_TOKEN"); v == "" {
		t.Fatal("NJALLA_API_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("NJALLA_TESTACC_DOMAIN"); v == "" {
		t.Fatal("NJALLA_TESTACC_DOMAIN must be set for acceptance tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
