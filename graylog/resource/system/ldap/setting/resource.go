package setting

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: create,
		Read:   read,
		Update: update,
		Delete: destroy,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// required
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Title of the LDAP backend",
			},
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "LDAP server URL (e.g., ldap://localhost:389 or ldaps://localhost:636)",
			},
			"authentication_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for binding to LDAP server",
			},
			"authentication_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for binding to LDAP server",
			},
			"user_search_base": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base DN for user searches",
			},
			"user_search_pattern": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pattern for searching users (e.g., uid={0})",
			},
			"user_name_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "LDAP attribute to use as the username",
			},
			"default_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the default role to assign to LDAP users",
			},

			// optional
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether this backend is enabled",
			},
			"user_id_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP attribute to use as the unique user ID",
			},
			"email_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP attribute to use as the email address",
			},
			"full_name_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP attribute to use as the full name",
			},
			"referral": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Referral mode (follow, ignore, etc.)",
			},
			"trust_fingerprint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trust fingerprint for the LDAP server certificate",
			},
			"start_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use StartTLS",
			},
			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection timeout in milliseconds",
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Request timeout in milliseconds",
			},
		},
	}
}
