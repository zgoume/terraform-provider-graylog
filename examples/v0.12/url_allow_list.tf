resource "graylog_url_allow_list" "example" {
  entries {
    id    = "alerta-webhook"
    title = "alerta"
    value = "https://alerta.cloud-temple.sec/api/webhooks/graylog"
    type  = "literal"
  }

  disabled = false
}
