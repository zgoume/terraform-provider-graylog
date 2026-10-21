# Resource: graylog_ldap_setting

## Example Usage

```hcl
resource "graylog_ldap_setting" "system" {
  title                     = "LDAP"
  server_url                = "ldap://localhost:389"
  authentication_username   = "cn=admin,dc=example,dc=com"
  authentication_password   = "secret"
  user_search_base          = "ou=users,dc=example,dc=com"
  user_search_pattern       = "uid={0}"
  user_name_attribute       = "uid"
  default_role              = "Reader"
  enabled                   = true
  user_id_attribute         = "entryUUID"
  email_attribute           = "mail"
  full_name_attribute       = "cn"
  start_tls                 = false
  connect_timeout           = 5000
  request_timeout           = 10000
}
```

## Argument Reference

* `title` - (Required) The data type is `string`. Title of the LDAP backend.
* `server_url` - (Required) The data type is `string`. LDAP server URL (e.g., `ldap://localhost:389` or `ldaps://localhost:636`).
* `authentication_username` - (Required) The data type is `string`. Username for binding to LDAP server.
* `authentication_password` - (Required) The data type is `string`. Password for binding to LDAP server. Sensitive field.
* `user_search_base` - (Required) The data type is `string`. Base DN for user searches.
* `user_search_pattern` - (Required) The data type is `string`. Pattern for searching users (e.g., `uid={0}`).
* `user_name_attribute` - (Required) The data type is `string`. LDAP attribute to use as the username.
* `default_role` - (Required) The data type is `string`. Name of the default role to assign to LDAP users.

* `enabled` - (Optional) The data type is `bool`. Whether this backend is enabled. Default: `true`.
* `user_id_attribute` - (Optional) The data type is `string`. LDAP attribute to use as the unique user ID.
* `email_attribute` - (Optional) The data type is `string`. LDAP attribute to use as the email address.
* `full_name_attribute` - (Optional) The data type is `string`. LDAP attribute to use as the full name.
* `start_tls` - (Optional) The data type is `bool`. Whether to use StartTLS.
* `connect_timeout` - (Optional) The data type is `int`. Connection timeout in milliseconds.
* `request_timeout` - (Optional) The data type is `int`. Request timeout in milliseconds.

## Attributes Reference

No additional attributes are exported.

## Import

LDAP settings can be imported using any id (the resource uses a fixed internal id):

```console
$ terraform import graylog_ldap_setting.foo bar
```

Note: When importing, you must specify the `title` in your Terraform config to match the existing backend.
