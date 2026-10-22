package setting

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/client"
)

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	roleUUID, err := findRoleUUID(ctx, cl.LDAPSetting.Client, d.Get("default_role").(string))
	if err != nil {
		return fmt.Errorf("failed to find role UUID for default role '%s': %w", d.Get("default_role").(string), err)
	}

	requestBody := map[string]interface{}{
		"title":         d.Get("title"),
		"default_roles": []string{roleUUID},
		"config":        buildLdapConfig(d),
	}

	if _, err := cl.LDAPSetting.CreateBackend(ctx, requestBody); err != nil {
		return fmt.Errorf("failed to create a LDAP backend: %w", err)
	}

	backends, _, err := cl.LDAPSetting.ListBackends(ctx)
	if err != nil {
		return fmt.Errorf("failed to list LDAP backends after creation: %w", err)
	}

	title := d.Get("title").(string)
	for _, b := range backends {
		if b["title"] == title {
			if backendID, ok := b["id"].(string); ok {
				if _, err := cl.LDAPSetting.SetActiveBackend(ctx, backendID); err != nil {
					return fmt.Errorf("failed to activate LDAP backend: %w", err)
				}
				break
			}
		}
	}

	d.SetId(ldapSettingID)
	return read(d, m)
}

func buildLdapConfig(d *schema.ResourceData) map[string]interface{} {
	scheme, host, port := parseServerURL(d.Get("server_url").(string))

	config := map[string]interface{}{
		"servers": []map[string]interface{}{
			{
				"host": host,
				"port": port,
			},
		},
		"type":                 "ldap",
		"system_user_dn":       d.Get("authentication_username").(string),
		"system_user_password": d.Get("authentication_password").(string),
	}

	if scheme == "ldaps" {
		config["transport_security"] = "tls"
	} else if d.Get("start_tls").(bool) {
		config["transport_security"] = "start_tls"
	} else {
		config["transport_security"] = "none"
	}

	if v := d.Get("user_id_attribute").(string); v != "" {
		config["user_unique_id_attribute"] = v
	}
	if v := d.Get("email_attribute").(string); v != "" {
		config["email_attributes"] = []string{v}
	}
	if v := d.Get("full_name_attribute").(string); v != "" {
		config["user_full_name_attribute"] = v
	}
	if v := d.Get("user_name_attribute").(string); v != "" {
		config["user_name_attribute"] = v
	}
	if v := d.Get("user_search_base").(string); v != "" {
		config["user_search_base"] = v
	}
	if v := d.Get("user_search_pattern").(string); v != "" {
		config["user_search_pattern"] = v
	}

	return config
}
