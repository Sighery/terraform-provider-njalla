package njalla

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Sighery/gonjalla"
)

func resourceRecordNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordNSCreate,
		Read:   resourceRecordNSRead,
		Update: resourceRecordNSUpdate,
		Delete: resourceRecordNSDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice(gonjalla.ValidTTL),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: resourceRecordNSImport,
		},
	}
}

func resourceRecordNSCreate(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("Adding record failed: %s", err.Error())
	}

	d.SetId(fmt.Sprint(saved.ID))

	return resourceRecordNSRead(d, m)

}

func resourceRecordNSRead(d *schema.ResourceData, m interface{}) error {
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

func resourceRecordNSUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	domain := d.Get("domain").(string)
	id, _ := strconv.Atoi(d.Id())

	updateRecord := gonjalla.Record{
		ID:      id,
		Name:    d.Get("name").(string),
		Type:    "NS",
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

	return resourceRecordNSRead(d, m)
}

func resourceRecordNSDelete(d *schema.ResourceData, m interface{}) error {
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

func resourceRecordNSImport(
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
