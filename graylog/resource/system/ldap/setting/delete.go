package setting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/client"
)

func destroy(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	// List all backends to find the one matching the title
	backends, _, err := cl.LDAPSetting.ListBackends(ctx)
	if err != nil {
		return fmt.Errorf("failed to list LDAP backends: %w", err)
	}

	title := d.Get("title").(string)
	for _, b := range backends {
		if b["title"] == title {
			backendID := b["id"].(string)
			if _, err := cl.LDAPSetting.DeleteBackend(ctx, backendID); err != nil {
				return fmt.Errorf("failed to delete a LDAP backend: %w", err)
			}
			return nil
		}
	}

	return nil
}
