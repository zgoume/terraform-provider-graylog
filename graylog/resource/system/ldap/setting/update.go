package setting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/client"
)

func update(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	backends, _, err := cl.LDAPSetting.ListBackends(ctx)
	if err != nil {
		return fmt.Errorf("failed to list LDAP backends: %w", err)
	}

	title := d.Get("title").(string)
	var backendID string
	for _, b := range backends {
		if b["title"] == title {
			backendID = b["id"].(string)
			break
		}
	}

	if backendID == "" {
		return fmt.Errorf("LDAP backend with title '%s' not found", title)
	}

	if _, err := cl.LDAPSetting.DeleteBackend(ctx, backendID); err != nil {
		return fmt.Errorf("failed to delete existing LDAP backend: %w", err)
	}

	return create(d, m)
}
