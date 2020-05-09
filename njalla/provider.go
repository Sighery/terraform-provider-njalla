package njalla

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider for Njalla resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Njalla API token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"njalla_record_txt": resourceRecordTXT(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("api_token").(string)

	config := Config{
		Token: token,
	}

	return &config, nil
}
