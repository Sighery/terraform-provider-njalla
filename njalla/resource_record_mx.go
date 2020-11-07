package njalla

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordMX() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordMXCreate,
		ReadContext:   resourceRecordMXRead,
		UpdateContext: resourceRecordMXUpdate,
		DeleteContext: resourceRecordMXDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the domain this record will be applied to.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: func() (interface{}, error) {
					return "@", nil
				},
				Description: "Name for the record.",
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "TTL for the record.",
				ValidateFunc: validation.IntInSlice(gonjalla.ValidTTL),
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Priority for the record.",
				ValidateFunc: validation.IntInSlice(gonjalla.ValidPriority),
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content for the record.",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordMXImport,
		},
	}
}

func resourceRecordMXCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	priority := d.Get("priority").(int)

	record := gonjalla.Record{
		Type:     "MX",
		Name:     d.Get("name").(string),
		Content:  d.Get("content").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &priority,
	}

	saved, err := gonjalla.AddRecord(config.Token, domain, record)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(saved.ID))

	return resourceRecordMXRead(ctx, d, m)

}

func resourceRecordMXRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	var diags diag.Diagnostics

	records, err := gonjalla.ListRecords(config.Token, domain)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, record := range records {
		if id == record.ID {
			d.Set("name", record.Name)
			d.Set("ttl", record.TTL)
			d.Set("priority", *record.Priority)
			d.Set("content", record.Content)

			return diags
		}
	}

	d.SetId("")
	return diags
}

func resourceRecordMXUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())
	priority := d.Get("priority").(int)

	updateRecord := gonjalla.Record{
		ID:       id,
		Name:     d.Get("name").(string),
		Type:     "MX",
		Content:  d.Get("content").(string),
		TTL:      d.Get("ttl").(int),
		Priority: &priority,
	}

	err := gonjalla.EditRecord(config.Token, domain, updateRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordMXRead(ctx, d, m)
}

func resourceRecordMXDelete(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	err := gonjalla.RemoveRecord(config.Token, domain, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceRecordMXImport(
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
			d.SetId(fmt.Sprintf("%d", id))
			d.Set("domain", domain)
			d.Set("name", record.Name)
			d.Set("ttl", record.TTL)
			d.Set("priority", *record.Priority)
			d.Set("content", record.Content)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find record %d for domain %s", id, domain)
}
