package njalla

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordNS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordNSCreate,
		ReadContext:   resourceRecordNSRead,
		UpdateContext: resourceRecordNSUpdate,
		DeleteContext: resourceRecordNSDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the domain this record will be applied to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for the record.",
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "TTL for the record.",
				ValidateFunc: validation.IntInSlice(gonjalla.ValidTTL),
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content for the record.",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordNSImport,
		},
	}
}

func resourceRecordNSCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	record := gonjalla.Record{
		Type:    "NS",
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	saved, err := gonjalla.AddRecord(config.Token, domain, record)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(saved.ID)

	return resourceRecordNSRead(ctx, d, m)

}

func resourceRecordNSRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	var diags diag.Diagnostics

	records, err := gonjalla.ListRecords(config.Token, domain)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, record := range records {
		if d.Id() == record.ID {
			d.Set("name", record.Name)
			d.Set("ttl", record.TTL)
			d.Set("content", record.Content)

			return diags
		}
	}

	d.SetId("")
	return diags
}

func resourceRecordNSUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	updateRecord := gonjalla.Record{
		ID:      d.Id(),
		Name:    d.Get("name").(string),
		Type:    "NS",
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	err := gonjalla.EditRecord(config.Token, domain, updateRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordNSRead(ctx, d, m)
}

func resourceRecordNSDelete(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	err := gonjalla.RemoveRecord(config.Token, domain, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceRecordNSImport(
	ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	domain, id, err := parseImportID(d.Id())
	if err != nil {
		return nil, err
	}

	config := m.(*Config)

	records, err := gonjalla.ListRecords(config.Token, domain)
	if err != nil {
		return nil, fmt.Errorf(
			"Reading records for domain %s failed: %s", domain, err.Error(),
		)
	}

	for _, record := range records {
		if id == record.ID {
			d.SetId(id)
			d.Set("domain", domain)
			d.Set("name", record.Name)
			d.Set("ttl", record.TTL)
			d.Set("content", record.Content)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find record %s for domain %s", id, domain)
}
