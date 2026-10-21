# graylog_url_allow_list

Manages the Graylog URL Allow List. This is a singleton resource - there is only one URL Allow List per Graylog instance.

## Example Usage

```hcl
resource "graylog_url_allow_list" "example" {
  entries {
    id    = "076419b6-8b36-4be6-b774-b3d8035640f1"
    title = "alerta"
    value = "https://alerta.cloud-temple.sec/api/webhooks/graylog"
    type  = "literal"
  }

  entries {
    id    = "102626ba-b749-411d-9d36-be2cb93572b0"
    title = "example-regex"
    value = "https://example\\.com/.*"
    type  = "regex"
  }

  disabled = false
}
```

## Argument Reference

The following arguments are supported:

* `entries` - (Required) The list of URL allow list entries. Configuration documentation below.
* `disabled` - (Optional) Whether the URL Allow List is disabled. Defaults to `false`.

**entries** supports the following:

* `id` - (Required) The unique identifier for the entry.
* `title` - (Required) A human-readable title for the entry.
* `value` - (Required) The URL or regex pattern to allow.
* `type` - (Required) The type of entry. Can be `literal` or `regex`.

## Attributes Reference

No additional attributes are exported.

## Import

URL Allow List can be imported using the id `url-allow-list`:

```bash
terraform import graylog_url_allow_list.example url-allow-list
```
