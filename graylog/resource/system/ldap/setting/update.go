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

	// List all backends to find the one matching the title
	backends, _, err := cl.LDAPSetting.ListBackends(ctx)
	if err != nil {
		return fmt.Errorf("failed to list LDAP backends: %w", err)
	}

	title := d.Get("title").(string)
	var backendID string
	var backend map[string]interface{}
	for _, b := range backends {
		if b["title"] == title {
			backend = b
			backendID = b["id"].(string)
			break
		}
	}

	if backend == nil {
		return fmt.Errorf("LDAP backend with title '%s' not found", title)
	}

	// Look up the default role UUID
	roleUUID, err := findRoleUUID(ctx, cl.LDAPSetting.Client, d.Get("default_role").(string))
	if err != nil {
		return fmt.Errorf("failed to find role '%s': %w", d.Get("default_role"), err)
	}

	// Build the request body
	requestBody := map[string]interface{}{
		"title":                   d.Get("title"),
		"type":                    "ldap",
		"enabled":                 d.Get("enabled"),
		"authentication_dn":       d.Get("authentication_username"),
		"authentication_password": d.Get("authentication_password"),
		"default_roles":           []string{roleUUID},
		"config":                  buildLdapConfig(d),
	}

	if _, err := cl.LDAPSetting.UpdateBackend(ctx, backendID, requestBody); err != nil {
		return fmt.Errorf("failed to update a LDAP backend: %w", err)
	}

	return read(d, m)
}
