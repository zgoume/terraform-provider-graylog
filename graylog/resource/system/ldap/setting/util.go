package setting

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/client/httpclient"
	"github.com/zahiar/terraform-provider-graylog/graylog/convert"
)

const (
	ldapSettingID = "ldap_setting_id"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, Resource(), flattenBackend(data)); err != nil {
		return err
	}

	d.SetId(ldapSettingID)
	return nil
}

func parseServerURL(serverURL string) (string, string, int) {
	scheme := "ldap"
	host := "localhost"
	port := 389

	parsed, err := url.Parse(serverURL)
	if err == nil {
		if parsed.Scheme == "ldaps" {
			scheme = "ldaps"
			port = 636
		} else if parsed.Scheme == "ldap" {
			scheme = "ldap"
			port = 389
		}
		if parsed.Hostname() != "" {
			host = parsed.Hostname()
		}
		if parsed.Port() != "" {
			fmt.Sscanf(parsed.Port(), "%d", &port)
		}
	}

	return scheme, host, port
}

func flattenBackend(backend map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if v, ok := backend["title"]; ok {
		result["title"] = v
	}
	if v, ok := backend["enabled"]; ok {
		result["enabled"] = v
	}

	if config, ok := backend["config"].(map[string]interface{}); ok {
		if servers, ok := config["servers"].([]map[string]interface{}); ok && len(servers) > 0 {
			server := servers[0]
			host := "localhost"
			port := 389
			if v, ok := server["host"]; ok {
				host = v.(string)
			}
			if v, ok := server["port"]; ok {
				fmt.Sscanf(fmt.Sprintf("%v", v), "%d", &port)
			}

			transportSecurity := "none"
			if v, ok := config["transport_security"]; ok {
				transportSecurity = v.(string)
			}
			if transportSecurity == "tls" {
				result["server_url"] = fmt.Sprintf("ldaps://%s:%d", host, port)
			} else {
				result["server_url"] = fmt.Sprintf("ldap://%s:%d", host, port)
			}
		}

		if v, ok := config["user_unique_id_attribute"]; ok {
			result["user_id_attribute"] = v
		}
		if v, ok := config["email_attributes"]; ok {
			if attrs, ok := v.([]interface{}); ok && len(attrs) > 0 {
				result["email_attribute"] = attrs[0]
			}
		}
		if v, ok := config["user_full_name_attribute"]; ok {
			result["full_name_attribute"] = v
		}
		if v, ok := config["user_name_attribute"]; ok {
			result["user_name_attribute"] = v
		}
		if v, ok := config["user_search_base"]; ok {
			result["user_search_base"] = v
		}
		if v, ok := config["user_search_pattern"]; ok {
			result["user_search_pattern"] = v
		}
		if v, ok := config["connect_timeout"]; ok {
			result["connect_timeout"] = v
		}
		if v, ok := config["request_timeout"]; ok {
			result["request_timeout"] = v
		}
		if v, ok := config["referral"]; ok {
			result["referral"] = v
		}
		if v, ok := config["trust_fingerprint"]; ok {
			result["trust_fingerprint"] = v
		}
	}

	if v, ok := backend["authentication_dn"]; ok {
		result["authentication_username"] = v
	}

	if roleIDs, ok := backend["default_roles"].([]interface{}); ok && len(roleIDs) > 0 {
		if roleMap, ok := backend["_roleMap"].(map[string]string); ok && len(roleMap) > 0 {
			roleID := roleIDs[0].(string)
			if roleName, ok := roleMap[roleID]; ok {
				result["default_role"] = roleName
			} else {
				result["default_role"] = roleID
			}
		} else {
			result["default_role"] = roleIDs[0].(string)
		}
	}

	delete(result, "_roleMap")

	return result
}

func findRoleUUID(ctx context.Context, cl httpclient.Client, roleName string) (string, error) {
	type Role struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type RolesResponse struct {
		Roles []Role `json:"roles"`
	}

	body := RolesResponse{}
	resp, err := cl.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/authz/roles?page=1&per_page=100",
		ResponseBody: &body,
	})
	if err != nil {
		if httpErr, ok := err.(interface {
			StatusCode() int
			BodyByte() []byte
		}); ok {
			return "", fmt.Errorf("failed to fetch roles (HTTP %d): %s: %w", httpErr.StatusCode(), string(httpErr.BodyByte()), err)
		}
		return "", err
	}
	_ = resp

	for _, role := range body.Roles {
		if role.Name == roleName {
			return role.ID, nil
		}
	}

	return "", fmt.Errorf("role '%s' not found", roleName)
}

func fetchRoleMap(ctx context.Context, cl httpclient.Client) (map[string]string, error) {
	type Role struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type RolesResponse struct {
		Roles []Role `json:"roles"`
	}

	body := RolesResponse{}
	resp, err := cl.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/authz/roles?page=1&per_page=100",
		ResponseBody: &body,
	})
	if err != nil {
		if httpErr, ok := err.(interface {
			StatusCode() int
			BodyByte() []byte
		}); ok {
			return nil, fmt.Errorf("failed to fetch roles (HTTP %d): %s: %w", httpErr.StatusCode(), string(httpErr.BodyByte()), err)
		}
		return nil, err
	}
	_ = resp

	roleMap := make(map[string]string)
	for _, role := range body.Roles {
		roleMap[role.ID] = role.Name
	}

	return roleMap, nil
}
