package njalla

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Sighery/gonjalla"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:		schema.TypeString,
				Required:	true,
				Description:	"Name for the server",
			},
			"instance_type": {
				Type:		schema.TypeString,
				Required:	true,
				Description:	"Instance type for the server",
			},
			"os": {
				Type:		schema.TypeString,
				Required:	true,
				Description:	"OS type for the server",
			},
			"public_key": {
				Type:		schema.TypeString,
				Required:	true,
				Description:	"Public key material for this server",
			},
			"months": {
				Type:		schema.TypeInt,
				Required:	true,
				Description:	"Number of months to buy the server for",
			},
			"public_ip": {
				Type:		schema.TypeString,
				Computed:	true,
				Description:	"Public IPv4 address of this server",
			},

		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImport,
		},
	}
}

func resourceServerCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	name := d.Get("name").(string)
	instanceType := d.Get("instance_type").(string)
	os := d.Get("os").(string)
	publicKey := d.Get("public_key").(string)
	months := d.Get("months").(int)

	server, err := gonjalla.AddServer(config.Token, name, instanceType, os, publicKey, months)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServerRead(ctx, d, m)
}

func resourceServerRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	servers, err := gonjalla.ListServers(config.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	for true {
		for _, server := range servers {
			if d.Id() == server.ID {
				if len(server.Ips) == 0 {
					time.Sleep(5 * time.Second)
					continue
				}

				d.Set("instance_type", server.Type)
				d.Set("os", server.Os)
				d.Set("public_key", server.SSHKey)
				d.Set("public_ip", server.Ips[0])

				return diags
			}
		}
	}

	d.SetId("")
	return diags
}

func resourceServerUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	id := d.Id()
	os := d.Get("os").(string)
	publicKey := d.Get("public_key").(string)
	instanceType := d.Get("instance_type").(string)

	_, err := gonjalla.ResetServer(config.Token, id, os, publicKey, instanceType)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServerRead(ctx, d, m)
}

func resourceServerDelete(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	config := m.(*Config)

	_, err := gonjalla.RemoveServer(config.Token, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func resourceServerImport(
	ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {

	id := d.Id()

	config := m.(*Config)

	servers, err := gonjalla.ListServers(config.Token)
	if err != nil {
		return nil, fmt.Errorf(
			"Listing servers failed: %s", err.Error(),
		)
	}

	for _, server := range servers {
		if id == server.ID {
			d.SetId(id)
			d.Set("name", server.Name)
			d.Set("instance_type", server.Type)
			d.Set("os", server.Os)
			d.Set("public_key", server.SSHKey)
			d.Set("public_ip", server.Ips[0])

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find server with id %s", id)
}
