# njalla_server Resource

Njalla Linux server/vps

## Example Usage

```hcl
resource "njalla_server" "teamserver" {
  name      = "teamserver"
  instance_type   = "njalla1"
  os        = "ubuntu2004"
  public_key    = tls_private_key.this.public_key_openssh
  months      = 1
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the server.
* `instance_type` - (Required) The njalla VPS instance type. 
    - "njalla1" = VPS 15 (1 core, 1.5 GB RAM, 15 GB Disk, 1.5 TB traffic)
    - "njalla2" = VPS 30 (2 cores, 3 GB RAM, 30 GB Disk, 3 TB traffic)
    - "njalla3" = VPS 45 (3 cores, 4.5 GB RAM, 45 GB Disk, 4.5 TB traffic)
* `os` - (Required) The operating system for the VPS
    - "debian11" = Debian 11
    - "debian12" = Debian 12
    - "ubuntu1804" = Ubuntu 18.04
    - "ubuntu2004" = Ubuntu 20.04
    - "ubuntu2204" = Ubuntu 22.04
    - "fedora38" = Fedora 38
    - "centos8stream" = CentOS 8 Stream
    - "rockylinux9" = Rocky Linux 9
    - "alpine318" = Alpine Linux 3.18
    - "archlinux" = Arch Linux
* `public_key` - (Required) The ssh public key material string for remote access.
* `months` - (Required) The number of months to rent the server for.

## Attributes Reference

* `id` - Njalla ID for this server.
* `name` - The name for this server.
* `instance_type` - The instance type for this server.
* `os` - The os type for this server.
* `public_key` - The public key material for this server.
* `public_ip` - The public IPv4 address for this server.

[gonjalla variable ValidTTL]: https://pkg.go.dev/github.com/Sighery/gonjalla?tab=doc#pkg-variables