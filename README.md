# Unofficial Terraform Njalla Provider

[Njalla][] is a privacy-oriented domain name registration service. Recently
they released their [official API][Njalla API]. Following the release of their
official API, I created an (for now extremely limited)
[unofficial Golang package][gonjalla package] for their API called `gonjalla`.

This repository is the unofficial Terraform provider for the Njalla API, using
the `gonjalla` package.

---

## Installing

Currently this provider isn't available in any distribution's repositories.

### Arch Linux (and derivatives)

Even though this package isn't in the official Arch Linux repositories (nor
the AUR), I maintain my own PKGBUILD and the compiled packages are in my Arch
Linux repository. More information in the
[terraform-provider-njalla-pkgbuild repository][terraform-provider-njalla-pkgbuild].

### Source

Make sure you have [Go][Golang] installed.

```bash
go build -o terraform-provider-njalla_v0.1.0
# Third-party Terraform plugins must be installed to this location
mkdir -p ~/.terraform.d/plugins
mv terraform-provider-njalla_v0.1.0 ~/.terraform.d/plugins/.
```

### Releases

There's a [Github Action set up to handle releases][Action Release] on tag
pushes. This Action currently builds and publishes the following binaries:

* `linux_amd64`

If you happen to run a system the Action currently builds for, you can go into
the [Releases tab][terraform-provider-njalla Releases] and download the
`.tar.gz` file for your system. Inside you'll find the built binary, along
with the LICENSE and README.

```bash
tar -xvf v0.1.0_linux_amd64.tar.gz
# Third-party Terraform plugins must be installed to this location
mkdir -p ~/.terraform.d/plugins
mv terraform-provider-njalla_v0.1.0 ~/.terraform.d/plugins/.
```

---

## Limitations

This provider only offers as much as is implemented in the `gonjalla` package.
That means that currently only the following resources are implemented:

* A record
* AAAA record
* TXT record
* MX record
* CNAME record
* CAA record
* PTR record
* NS record
* TLSA record

If you have a need for any other resource, please feel free to contribute to
both the `gonjalla` and this repository.

## Usage

### Setting up the provider

The Njalla API uses a token you can generate from the settings page. This
token can be set up in two ways:

#### Terraform configuration

On your `config.tf` (name just a convention) file:

```terraform
provider njalla {
  api_token = "api-token-here"
}
```

#### Environment variables

Or through the use of the environment variable `NJALLA_API_TOKEN`. If this
environment variable is set you don't need to set up the `njalla` provider in
your Terraform files, it will be automatically picked up when initialising the
code.

### Record A

#### Basic

```terraform
resource njalla_record_a example-a {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "138.201.81.199"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) IPv4 address for the record.

### Record AAAA

#### Basic

```terraform
resource njalla_record_aaaa example-aaaa {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "2a01:4f8:172:1d86::1"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) IPv6 address for the record.

### Record TXT

#### Basic

```terraform
resource njalla_record_txt example-txt {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-content"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record.

### Record MX

#### Basic

```terraform
resource njalla_record_mx example-mx {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  priority = 10
  content = "example-content"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `priority` - (Required) Priority for the record. Value must be one of
  [gonjalla `ValidPriority`][gonjalla variable ValidPriority].
* `content` - (Required) Content for the record.

### Record CNAME

#### Basic

```terraform
resource njalla_record_cname example-cname {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-website.com"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record.

### Record CAA

#### Basic

```terraform
resource njalla_record_caa example-caa {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-content"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record. Value must follow the
  [RFC 8659][]'s syntax from point 4.

### Record PTR

#### Basic

```terraform
resource njalla_record_ptr example-ptr {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-content"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record.

### Record NS

#### Basic

```terraform
resource njalla_record_ns example-ns {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "example-content"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Required) Name for the record.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record.

### Record TLSA

#### Basic

```terraform
resource njalla_record_tlsa example-tlsa {
  domain = "example.com"
  name = "example-name"
  ttl = 10800
  content = "0 0 1 d2abde240d7cd3ee6b4b28c54df034b9"
}
```

#### Argument Reference

* `domain` - (Required) Specifies the domain this record will be applied to.
  Changing this forces a new resource to be created.
* `name` - (Optional) Name for the record. Default is `@`.
* `ttl` - (Required) TTL for the record. Value must be one of
  [gonjalla `ValidTTL`][gonjalla variable ValidTTL].
* `content` - (Required) Content for the record. Value must follow
  [RFC 6698][]'s syntax from sections 2 and 7.

### Importing existing resources

Currently all the available resources implement the import functionality.
[More information on Terraform's import][Terraform import].

Imagine the following example: You have set up DNS records manually through
the Njalla web interface, and you're now trying to manage them through
Terraform. To do so, first you need to declare the resource in your Terraform
code:

```terraform
resource njalla_record_txt example-import {}
```

As you may have noticed, the declared resource is empty. Your linter might be
yelling at you now. Ignore it for now and leave it as is.

Now, to import the DNS record you want to attach to this new resource, you'd
execute the following:

```bash
# Base command
terraform import address domain:id
# Example command
terraform import njalla_record_txt.example-import example.com:12345
```

Check the [Terraform import usage][terraform import] on how to figure out the
`address` positional argument. It will depend on your Terraform declaration
and where have you defined your resources.

The `domain:id` bit is the important part, and specific to this provider.
Since records are attached to a given domain, to import a record into
Terraform we need both the Njalla ID of the record, and the domain it's
attached to.

#### Import only updates the State file

**Please note**: After executing this command, you might notice your Terraform
code hasn't changed in the slightest. This is a
[known limitation of Terraform][Terraform import state only limitation] as of
now. The State file will be updated to have your imported resource, but the
Terraform code will look just the same.

Fill in the resource manually now, you can use the data in the State file to
know what to type in the Terraform code. If you don't, on your next execution
Terraform will try to remove or update the resource.

Even with this limitation, this is still useful for letting Terraform know
that it can go ahead and manage an already existing resource, without having
to resort to manually deleting the resource and recreating it with Terraform
for it to know it can manage it now.

---

## Contributing

As mentioned previously, this provider plugin depends completely on the
[gonjalla package][]. If you wanted to add new Njalla resources to this
provider, chances are you'd first have to implement them in the `gonjalla`
package.

Assuming you've done that, and followed that package's contributing guides,
once adding new resources to the provider, here's how I do it.

### New resource

Add any new resources inside the `njalla` package. The file name must follow
this format: `resource_{type}`. In the case of our `njalla_record_txt`
resource, the file is then called [`resource_record_txt.go`][]. Take a look at
any of the existing resources, and how they're linked in [`provider.go`][].

You'll have to implement the basic `CRUD` operations, and if possible, **do
implement importing as well**.

### Acceptance tests

After adding your new resource (or before), **add acceptance tests**.
Take a look at the
[documentation][Terraform provider acceptance tests documentation] to learn
more. [This Medium article][Terraform provider acceptance tests article]
helped me greatly to understand the structure of acceptance tests and the
whole complex system around them.

Do take a look at existing acceptance tests for the implemented record types,
and copy however much is useful, editing when needed. When copying tests,
**remember to add tests for any new or specific functionality to the new
resource**. For instance: If your new resource has a new field that can only
take certain values, write a new acceptance test for that functionality
specific to that new record.

### Running acceptance tests

These tests **will deploy new infrastructure**. They might fail, and leave
detached/dangling infrastructure, especially during development if your tests
are not yet working properly. It's up to you to clean up afterwards if this is
the case. Any acceptance tests after the development/testing phase is over
should not ever leave dangling resources. Please do test extensively before
making a pull request.

In this repository
[I have GitHub Action set up to run acceptance and unit tests][Action Test].
This action makes use of Action Secrets with a given Njalla API token and a
test domain I've set up to run these acceptance tests.

Acceptance tests still use Golang's testing functionality, which can be run by
executing:

```bash
go test -v ./...
```

However, when doing this, you might notice that Terraform's acceptance tests
are simply being skipped. This is because Terraform's SDK requires the
environment variable `TF_ACC` to be set to `true` to run acceptance tests.

This provider requires another two environment variables set to run acceptance
tests:

* `NJALLA_API_TOKEN`: Njalla API token used to call the API during tests.
* `NJALLA_TESTACC_DOMAIN`: Njalla domain used during the tests.

```bash
export NJALLA_API_TOKEN="api-token-here"
export NJALLA_TESTACC_DOMAIN="testdomain.com"
TF_ACC=true go test -v ./...
```

[Njalla]: https://njal.la
[Njalla API]: https://njal.la/api/
[gonjalla package]: https://github.com/Sighery/gonjalla
[terraform-provider-njalla-pkgbuild]: https://github.com/Sighery/terraform-provider-njalla-pkgbuild
[Golang]: https://golang.org/
[installing Terraform plugins]: https://www.terraform.io/docs/plugins/basics.html#installing-plugins
[gonjalla variable ValidTTL]: https://pkg.go.dev/github.com/Sighery/gonjalla?tab=doc#pkg-variables
[gonjalla variable ValidPriority]: https://pkg.go.dev/github.com/Sighery/gonjalla?tab=doc#pkg-variables
[RFC 8659]: https://tools.ietf.org/html/rfc8659
[RFC 6698]: https://tools.ietf.org/html/rfc6698
[Terraform import]: https://www.terraform.io/docs/import/usage.html
[Terraform import state only limitation]: https://www.terraform.io/docs/import/index.html#currently-state-only
[`resource_record_txt.go`]: njalla/resource_record_txt.go
[`provider.go`]: njalla/provider.go
[Terraform provider acceptance tests documentation]: https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html
[Terraform provider acceptance tests article]: https://medium.com/spaceapetech/creating-a-terraform-provider-part-2-1346f89f082c
[Action Test]: .github/workflows/test.yml
[Action Release]: .github/workflows/release.yml
[terraform-provider-njalla releases]: https://github.com/Sighery/terraform-provider-njalla/releases
