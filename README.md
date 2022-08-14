# Unofficial Terraform Njalla Provider

[Njalla][] is a privacy-oriented domain name registration service. Recently
they released their [official API][Njalla API]. Following the release of their
official API, I created an (for now extremely limited)
[unofficial Golang package][gonjalla package] for their API called `gonjalla`.

This repository is the unofficial Terraform provider for the Njalla API, using
the `gonjalla` package.

---

## Installing

Starting from Terraform v0.13, there's now a
[registry for providers][Terraform providers registry], where 
[this provider gets uploaded to][Sighery/njalla registry]. To use in your
Terraform project:

```terraform
terraform {
  required_version = ">= 0.13"

  required_providers {
    njalla = {
      source  = "Sighery/njalla"
      version = "~> 0.10.0"
    }
  }
}
```

With this, Terraform will now take care of finding the relevant provider in
the registry, and download it. After that, it can be configured
(`provider njalla {}` block) and used throughout the project.


## Documentation

The documentation is
[rendered online in the Terraform Registry][rendered documentation], generated
from the files in the [`docs/`][] directory, where you can render the Markdown
files locally to read them as well.

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

### Releasing

There's a [Github Action set up to handle releases][Action Release] on tag
pushes. This action then makes use of [GoReleaser][] to cross-compile to
different platforms. GoReleaser is also used to create a checksums file, and
create a new draft GitHub release.

From there, I'll download the checksums file, sign it with the GPG key linked
to this provider, and upload the signature file back into the release before
publishing it. Once published, the Registry website picks up the new release
automatically.

All providers in the registry must have a linked GPG key, and all the releases
for that provider must contain a signature file of the checksum signed by that
configured GPG key.
[More information here][provider signing key documentation].

[Njalla]: https://njal.la
[Njalla API]: https://njal.la/api/
[gonjalla package]: https://github.com/Sighery/gonjalla
[Terraform providers registry]: https://registry.terraform.io/browse/providers
[Sighery/njalla registry]: https://registry.terraform.io/providers/Sighery/njalla
[rendered documentation]: https://registry.terraform.io/providers/Sighery/njalla/latest/docs
[`docs/`]: docs/
[`resource_record_txt.go`]: njalla/resource_record_txt.go
[`provider.go`]: njalla/provider.go
[Terraform provider acceptance tests documentation]: https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html
[Terraform provider acceptance tests article]: https://medium.com/spaceapetech/creating-a-terraform-provider-part-2-1346f89f082c
[Action Test]: .github/workflows/test.yml
[Action Release]: .github/workflows/release.yml
[terraform-provider-njalla releases]: https://github.com/Sighery/terraform-provider-njalla/releases
[GoReleaser]: https://goreleaser.com/
[provider signing key documentation]: https://www.terraform.io/docs/registry/providers/publishing.html#preparing-and-adding-a-signing-key
