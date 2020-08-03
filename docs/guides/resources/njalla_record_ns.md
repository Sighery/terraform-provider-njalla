# njalla_record_ns Resource

Njalla `NS` DNS record for a given domain.

## Example Usage

```hcl
resource njalla_record_ns example-ns {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-content"
}
```

## Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla's `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record.

~> **Note** Changing the `domain` attribute forces the existing resource to be
deleted from the previous domain, and created into the new domain.

## Attributes Reference

* `id` - Njalla ID for this record.

[gonjalla variable ValidTTL]: https://pkg.go.dev/github.com/Sighery/gonjalla?tab=doc#pkg-variables
