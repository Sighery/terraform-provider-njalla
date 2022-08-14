---
page_title: "Importing existing Record resources"
---

# Importing Record resources

-> **Note** Currently importing only updates the State file. This is subject
to change in the future according to Terraform's roadmap. But at the moment it
means that after importing a resource, you'll have to type out manually the
resource definition in your Terraform files. If you don't, Terraform will
assume you want to delete the newly imported resource (since its local
definition doesn't match its remote status).

You may have set up DNS records manually through the Njalla web interface, and
you're now trying to migrate to managing them through Terraform. Luckily
Terraform allows for importing existing resources. To do so, first you need to
declare the resource in your Terraform code:

```hcl
resource njalla_record_txt example-import {}
```

The declared resource is empty for now. We just need it as a placeholder to
let Terraform know to attach a given remote resource to this new local
resource.

We now need to run the `import` command to pull the remote resource and attach
it to the local one:

```sh
# Base command
$ terraform import address domain:id
# Example command
$ terraform import njalla_record_txt.example-import example.com:12345
```

Check the [Terraform import usage][Terraform import] on how to figure out the
`address` positional argument. It will depend on your Terraform declaration
and where have you defined your resources.

The `domain:id` bit is the important part, and specific to this provider.
Since records are attached to a given domain, to import a record into
Terraform we need both the Njalla ID of the record, and the domain it's
attached to.

The `id` bit is the ID the record has in Njalla's DB. This, as of writing, can
only be fetched from the `list-records` API call, which will contain the value
under the key `id` for each record.

[Terraform import]: https://www.terraform.io/docs/import/usage.html
