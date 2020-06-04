---
layout: "ibm"
page_title: "IBM : dedicated_host"
sidebar_current: "docs-ibm-resource-is-dedicated_host"
description: |-
  Manages IBM IS Dedicated Host.
---

# ibm\_is_dedicated_host

Provides a dedicated host resource. This allows dedicated host to be created, updated, and destroyed.


## Example Usage

```hcl

resource "ibm_is_dedicated_host" "test_dedicated_host" {
  name = "testdedicatedhost"
  instance_placement_enabled = "false"
  group = ibm_is_dedicated_host_group.testdedicatedhostgroup.id
  resource_group = data.ibm_resource_group.resource_group.id
  profile = "dh2-56x464"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The dedicated host name.
* `instance_placement_enabled` - (Required, string) The instance placement enablement attribute. If set to true, instances can be placed on this dedicated host 
* `group` - (Required, string) The dedicated host group id for this dedicated host. 
* `profile` - (Required, string) The profile name. 
* `resource_group` - (Required, string) - The resource group id of this dedicated host. 

## Attribute Reference

The following attributes are exported:

* `id` - The id of the instance.
* `memory` - Memory of the instance.
* `status` - Status of the instance.
* `vcpu` - A nested block describing the VCPU configuration of this instance.
Nested `vcpu` blocks have the following structure:
  * `architecture` - The architecture of the instance.
  * `count` - The number of VCPUs assigned to the instance.
* `gpu` - A nested block describing the gpu of this instance.
Nested `gpu` blocks have the following structure:
  * `cores` - The cores of the gpu.
  * `count` - Count of the gpu.
  * `manufacture` - Manufacture of the gpu.
  * `memory` - Memory of the gpu.
  * `model` - Model of the gpu.
* `primary_network_interface` - A nested block describing the primary network interface of this instance.
Nested `primary_network_interface` blocks have the following structure:
  * `id` - The id of the network interface.
  * `name` - The name of the network interface.
  * `subnet` -  ID of the subnet.
  * `security_groups` -  List of security groups.
  * `primary_ipv4_address` - The primary IPv4 address.
* `network_interfaces` - A nested block describing the additional network interface of this instance.
Nested `network_interfaces` blocks have the following structure:
  * `id` - The id of the network interface.
  * `name` - The name of the network interface.
  * `subnet` -  ID of the subnet.
  * `security_groups` -  List of security groups.
  * `primary_ipv4_address` - The primary IPv4 address.
* `boot_volume` - A nested block describing the boot volume.
Nested `boot_volume` blocks have the following structure:
  * `name` - The name of the boot volume.
  * `size` -  Capacity of the volume in GB.
  * `iops` -  Input/Output Operations Per Second for the volume.
  * `profile` - The profile of the volume.
  * `encryption` - The encryption of the boot volume.
* `volume_attachments` - A nested block describing the volume attachments.  
Nested `volume_attachments` block have the following structure:
  * `id` - The id of the volume attachment
  * `name` -  The name of the volume attachment
  * `volume_id` - The id of the volume attachment's volume
  * `volume_name` -  The name of the volume attachment's volume
  * `volume_crn` -  The CRN of the volume attachment's volume
* `resource_controller_url` - The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance.


## Import

ibm_is_dedicated_host can be imported using dedicated host ID, eg

```
$ terraform import ibm_is_instance.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
