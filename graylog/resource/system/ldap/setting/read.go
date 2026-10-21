package setting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/client"
	"github.com/zahiar/terraform-provider-graylog/graylog/util"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	backends, resp, err := cl.LDAPSetting.ListBackends(ctx)
	if err != nil {
		return util.HandleGetResourceError(
			d, resp, fmt.Errorf("failed to list LDAP backends: %w", err))
	}

	title := d.Get("title").(string)
	var backend map[string]interface{}
	for _, b := range backends {
		if b["title"] == title {
			backend = b
			break
		}
	}

	if backend == nil {
		d.SetId("")
		return nil
	}

	roleMap, err := fetchRoleMap(ctx, cl.LDAPSetting.Client)
	if err != nil {
		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	backend["_roleMap"] = roleMap

	if err := setDataToResourceData(d, backend); err != nil {
		return fmt.Errorf("failed to set data: %w", err)
	}

	d.SetId(ldapSettingID)
	return nil
}
