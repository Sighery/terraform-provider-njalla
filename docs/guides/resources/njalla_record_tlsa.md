# njalla_record_tlsa Resource

Njalla `TLSA` DNS record for a given domain.

## Example Usage

```hcl
resource njalla_record_tlsa example-tlsa {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "0 0 1 d2abde240d7cd3ee6b4b28c54df034b9"
}
```

## Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla's `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record. Value must follow
  [RFC 6698][]'s syntax from sections 2 and 7.

~> **Note** Changing the `domain` attribute forces the existing resource to be
deleted from the previous domain, and created into the new domain.

## Attributes Reference

* `id` - Njalla ID for this record.

[gonjalla variable ValidTTL]: https://pkg.go.dev/github.com/Sighery/gonjalla?tab=doc#pkg-variables
[RFC 6698]: https://tools.ietf.org/html/rfc6698
