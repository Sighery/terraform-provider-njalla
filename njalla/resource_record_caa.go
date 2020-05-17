package njalla

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordCAA() *schema.Resource {
	contentRegex := regexp.MustCompile(
		`^\d{1,3}\s+(?:issue|iodef|issuewild)\s+.+$`,
	)

	return &schema.Resource{
		Create: resourceRecordCAACreate,
		Read:   resourceRecordCAARead,
		Update: resourceRecordCAAUpdate,
		Delete: resourceRecordCAADelete,

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
				Type:     schema.TypeString,
				Required: true,
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
			State: resourceRecordCAAImport,
		},
	}
}

func resourceRecordCAACreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("Adding record failed: %s", err.Error())
	}

	d.SetId(fmt.Sprint(saved.ID))

	return resourceRecordCAARead(d, m)

}

func resourceRecordCAARead(d *schema.ResourceData, m interface{}) error {
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

func resourceRecordCAAUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	updateRecord := gonjalla.Record{
		ID:      id,
		Name:    d.Get("name").(string),
		Type:    "CAA",
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

	return resourceRecordCAARead(d, m)
}

func resourceRecordCAADelete(d *schema.ResourceData, m interface{}) error {
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

func resourceRecordCAAImport(
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
