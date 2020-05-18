package njalla

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordNAPTR() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordNAPTRCreate,
		Read:   resourceRecordNAPTRRead,
		Update: resourceRecordNAPTRUpdate,
		Delete: resourceRecordNAPTRDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: func() (interface{}, error) {
					return "@", nil
				},
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice(gonjalla.ValidTTL),
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNAPTRContent,
			},
		},

		Importer: &schema.ResourceImporter{
			State: resourceRecordNAPTRImport,
		},
	}
}

func resourceRecordNAPTRCreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("Adding record failed: %s", err.Error())
	}

	d.SetId(fmt.Sprint(saved.ID))

	return resourceRecordNAPTRRead(d, m)

}

func resourceRecordNAPTRRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	records, err := gonjalla.ListRecords(config.Token, domain)
	if err != nil {
		return fmt.Errorf(
			"Reading records for domain %s failed: %s", domain, err.Error(),
		)
	}

	for _, record := range records {
		if id == record.ID {
			d.Set("name", record.Name)
			d.Set("ttl", record.TTL)
			d.Set("content", record.Content)

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceRecordNAPTRUpdate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf(
			"Updating record %d for domain %s failed: %s",
			id, domain, err.Error(),
		)
	}

	return resourceRecordNAPTRRead(d, m)
}

func resourceRecordNAPTRDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	err := gonjalla.RemoveRecord(config.Token, domain, id)
	if err != nil {
		return fmt.Errorf(
			"Deleting record %d from domain %s failed: %s",
			id, domain, err.Error(),
		)
	}

	return nil
}

func resourceRecordNAPTRImport(
	d *schema.ResourceData, m interface{},
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
