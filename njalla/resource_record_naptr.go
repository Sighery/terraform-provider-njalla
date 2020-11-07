package njalla

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordNAPTR() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordNAPTRCreate,
		ReadContext:   resourceRecordNAPTRRead,
		UpdateContext: resourceRecordNAPTRUpdate,
		DeleteContext: resourceRecordNAPTRDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Content for the record.",
				ValidateFunc: validateNAPTRContent,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordNAPTRImport,
		},
	}
}

func resourceRecordNAPTRCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)

	record := gonjalla.Record{
		Type:    "NAPTR",
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	saved, err := gonjalla.AddRecord(config.Token, domain, record)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(saved.ID))

	return resourceRecordNAPTRRead(ctx, d, m)

}

func resourceRecordNAPTRRead(
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
			d.Set("content", record.Content)

			return diags
		}
	}

	d.SetId("")
	return diags
}

func resourceRecordNAPTRUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	updateRecord := gonjalla.Record{
		ID:      id,
		Name:    d.Get("name").(string),
		Type:    "NAPTR",
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	err := gonjalla.EditRecord(config.Token, domain, updateRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordNAPTRRead(ctx, d, m)
}

func resourceRecordNAPTRDelete(
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

func resourceRecordNAPTRImport(
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
			d.Set("content", record.Content)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find record %d for domain %s", id, domain)
}

// validateNAPTRContent will be the `ValidateFunc` used to check a given
// content for a NAPTR DNS record matches the specification. If you're up for
// some heavy reading, check RFC 2915 section 2:
// https://tools.ietf.org/html/rfc2915
func validateNAPTRContent(
	val interface{}, key string,
) (warns []string, errs []error) {
	v := val.(string)
	values := strings.Split(v, " ")

	rfc := "Check RFC 2915 section 2"

	if len(values) < 6 {
		msg := fmt.Errorf(
			"expected 6+ arguments, got: %d. %s",
			len(values), rfc,
		)
		errs = append(errs, msg)
		return
	}

	_, err := strconv.Atoi(values[0])
	if err != nil {
		msg := fmt.Errorf(
			"expected Order field to be int, got: %s. %s",
			values[0], rfc,
		)
		errs = append(errs, msg)
		return
	}

	_, err = strconv.Atoi(values[1])
	if err != nil {
		msg := fmt.Errorf(
			"expected Preference field to be int, got: %s. %s",
			values[1], rfc,
		)
		errs = append(errs, msg)
		return
	}
	return
}
