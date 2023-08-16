package njalla

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"njalla_server":       resourceServer(),
			"njalla_domain":       resourceDomain(),
			"njalla_record_txt":   resourceRecordTXT(),
			"njalla_record_a":     resourceRecordA(),
			"njalla_record_aaaa":  resourceRecordAAAA(),
			"njalla_record_mx":    resourceRecordMX(),
			"njalla_record_cname": resourceRecordCNAME(),
			"njalla_record_caa":   resourceRecordCAA(),
			"njalla_record_ptr":   resourceRecordPTR(),
			"njalla_record_ns":    resourceRecordNS(),
			"njalla_record_tlsa":  resourceRecordTLSA(),
			"njalla_record_naptr": resourceRecordNAPTR(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := d.GetOk("api_token"); ok {
		token := v.(string)
		config := Config{
			Token: token,
		}

		return &config, diags
	}

	// Reaching here means the token wasn't given through the Terraform
	// config NOR environment variable (`DefaultFunc`).
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to setup Njalla provider",
		Detail:   "Missing required API token for provider Njalla",
	})
	return nil, diags
}
