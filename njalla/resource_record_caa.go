package njalla

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordCAA() *schema.Resource {
	contentRegex := regexp.MustCompile(
		`^\d{1,3}\s+(?:issue|iodef|issuewild)\s+.+$`,
	)

	return &schema.Resource{
		CreateContext: resourceRecordCAACreate,
		ReadContext:   resourceRecordCAARead,
		UpdateContext: resourceRecordCAAUpdate,
		DeleteContext: resourceRecordCAADelete,

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
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content for the record.",
				ValidateFunc: validation.All(
					validation.StringMatch(
						contentRegex,
						"value must follow RFC 8659: point 4 for syntax",
					),
					func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)
						r := regexp.MustCompile(`^(\d{1,3})\s+`)
						matches := r.FindStringSubmatch(v)

						if matches == nil || len(matches) < 2 {
							missingFlag := fmt.Errorf(
								"no flag found: RFC 8659 point 4.1.1",
							)
							errs = append(errs, missingFlag)
							return
						}

						flag, err := strconv.Atoi(matches[1])
						if err != nil {
							invalidFlag := fmt.Errorf(
								"flag is not int: RFC 8659 point 4.1.1",
							)
							errs = append(errs, invalidFlag)
							return
						}

						if 0 > flag || flag > 255 {
							invalidFlag := fmt.Errorf(
								"flag must be between 0 and 255: RFC 8659 4.1.1",
							)
							errs = append(errs, invalidFlag)
							return
						}

						return
					},
				),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordCAAImport,
		},
	}
}

func resourceRecordCAACreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	record := gonjalla.Record{
		Type:    "CAA",
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	saved, err := gonjalla.AddRecord(config.Token, domain, record)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(saved.ID)

	return resourceRecordCAARead(ctx, d, m)

}

func resourceRecordCAARead(
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

func resourceRecordCAAUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	updateRecord := gonjalla.Record{
		ID:      d.Id(),
		Name:    d.Get("name").(string),
		Type:    "CAA",
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	err := gonjalla.EditRecord(config.Token, domain, updateRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordCAARead(ctx, d, m)
}

func resourceRecordCAADelete(
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

func resourceRecordCAAImport(
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
