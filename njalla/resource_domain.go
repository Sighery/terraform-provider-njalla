package njalla

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Sighery/gonjalla"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: domainCreate,
		ReadContext:   domainRead,
		UpdateContext:   domainUpdate,
		DeleteContext: domainDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:		schema.TypeString,
				Required:	true,
				Description:	"Name of the domain",
			},
			"years": {
				Type:		schema.TypeInt,
				Required:	true,
				Description:	"Number of months to buy the server for",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: domainImport,
		},
	}
}

func domainCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	name := d.Get("name").(string)
	years := d.Get("years").(int)

	err := gonjalla.RegisterDomain(config.Token, name, years)
	if err != nil {
		return diag.FromErr(err)
	}

	return domainRead(ctx, d, m)
}

func domainRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	domains, err := gonjalla.ListDomains(config.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, domain := range domains {
		if d.Id() == domain.Name {
			d.Set("name", domain.Name)
		}

		return diags
	}

	d.SetId("")
	return diags
}

func domainUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {

	return domainRead(ctx, d, m)
}


func domainDelete(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {

	// Cannot delete domain through API.  Just delete in internal state instead.
	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func domainImport(
	ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {

	name := d.Id()

	config := m.(*Config)

	domains, err := gonjalla.ListDomains(config.Token)
	if err != nil {
		return nil, fmt.Errorf(
			"Listing domains failed: %s", err.Error(),
		)
	}

	for _, domain := range domains {
		if name == domain.Name {
			d.SetId(domain.Name)
			d.Set("name", domain.Name)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find domain with id %s", name)
}
