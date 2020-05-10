package njalla

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider for Njalla resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NJALLA_API_TOKEN", nil),
				Description: "Njalla API token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"njalla_record_txt": resourceRecordTXT(),
			"njalla_record_a": resourceRecordA(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	if v, ok := d.GetOk("api_token"); ok {
		token := v.(string)
		config := Config{
			Token: token,
		}

		return &config, nil
	}

	// Reaching here means the token wasn't given through the Terraform
	// config NOR environment variable (`DefaultFunc`).
	return nil, fmt.Errorf("Missing required API token for provider Njalla")
}
