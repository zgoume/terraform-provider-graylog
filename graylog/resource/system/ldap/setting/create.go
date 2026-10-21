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

	if _, err := cl.LDAPSetting.CreateBackend(ctx, requestBody); err != nil {
		return fmt.Errorf("failed to create a LDAP backend: %w", err)
	}

	d.SetId(ldapSettingID)
	return read(d, m)
}

// buildLdapConfig builds the config object for the LDAP backend
func buildLdapConfig(d *schema.ResourceData) map[string]interface{} {
	scheme, host, port := parseServerURL(d.Get("server_url").(string))

	config := map[string]interface{}{
		"servers": []map[string]interface{}{
			{
				"host": host,
				"port": port,
			},
		},
		"type": "ldap",
	}

	// Set transport_security based on scheme and start_tls
	if scheme == "ldaps" {
		config["transport_security"] = "tls"
	} else if d.Get("start_tls").(bool) {
		config["transport_security"] = "start_tls"
	} else {
		config["transport_security"] = "none"
	}

	// Set optional config fields
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
	if v := d.Get("connect_timeout").(int); v > 0 {
		config["connect_timeout"] = v
	}
	if v := d.Get("request_timeout").(int); v > 0 {
		config["request_timeout"] = v
	}
	if v := d.Get("referral").(string); v != "" {
		config["referral"] = v
	}
	if v := d.Get("trust_fingerprint").(string); v != "" {
		config["trust_fingerprint"] = v
	}

	return config
}
