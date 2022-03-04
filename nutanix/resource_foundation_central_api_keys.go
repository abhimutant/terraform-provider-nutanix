package nutanix

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/fc"
)

func resourceFoundationCentralAPIKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoundationCentralAPIKeysCreate,
		ReadContext:   resourceFoundationCentralAPIKeysRead,
		DeleteContext: resourceFoundationCentralAPIKeysDelete,
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFoundationCentralAPIKeysCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FC
	req := &fc.CreateAPIKeysInput{}

	alias, ok := d.GetOk("alias")
	if ok {
		req.Alias = alias.(string)
	}

	resp, err := conn.CreateAPIKey(req)
	if err != nil {
		return diag.Errorf("error creating API Keys with alias %s: %+v", (req.Alias), err)
	}

	d.SetId(resp.KeyUUID)
	return resourceFoundationCentralAPIKeysRead(ctx, d, meta)
}

func resourceFoundationCentralAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FC
	resp, err := conn.GetAPIKey(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("created_timestamp", resp.CreatedTimestamp)
	d.Set("key_uuid", resp.KeyUUID)
	d.Set("api_key", resp.ApiKey)
	d.Set("current_time", resp.CurrentTime)
	d.Set("alias", resp.Alias)

	return nil
}

func resourceFoundationCentralAPIKeysDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
