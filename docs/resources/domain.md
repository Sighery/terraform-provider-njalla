# njalla_record_a Resource

Njalla `A` DNS record for a given domain.

# WARNING !!!
There is not a way to delete a domain from njalla, so when a domain is deleted it is simply removed from the terraform state file. Furthermore, auto-renew is disabled by default.

## Example Usage

```hcl
resource "njalla_domain" "this" {
  name    = "example.com"
  years   = 1
}
```

## Argument Reference

* `name` - (Required) Specifies the domain name.
* `years` - (Required) Specifies the number of years to register the domain for.

## Attributes Reference

* `id` - Njalla ID for this domain.  This is the domain name.
* `name` - The domain name.