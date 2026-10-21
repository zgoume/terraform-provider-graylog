resource "graylog_ldap_setting" "system" {
  title                   = "LDAP"
  server_url              = "ldap://localhost:389"
  authentication_username = "cn=admin,dc=example,dc=com"
  authentication_password = "secret"
  user_search_base        = "ou=users,dc=example,dc=com"
  user_search_pattern     = "uid={0}"
  user_name_attribute     = "uid"
  default_role            = "Reader"
  enabled                 = true
  user_id_attribute       = "entryUUID"
  email_attribute         = "mail"
  full_name_attribute     = "cn"
  start_tls               = false
  connect_timeout         = 5000
  request_timeout         = 10000
}
