---
layout: "ibm"
page_title: "IBM : dedicated_hosts"
sidebar_current: "docs-ibm-datasources-is-dedicated_hosts"
description: |-
  List Dedicated hosts.
---

# ibm\_is_dedicated_hosts

Import the details of an existing IBM dedicated hosts as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```hcl

resource "ibm_is_dedicated_host" "test_dedicated_host" {
  name = "testdedicatedhost"
  instance_placement_enabled = "false"
  group = ibm_is_dedicated_host_group.test_dedicated_host_group.id
  resource_group = data.ibm_resource_group.rg.id
  profile = "dh2-56x464"
}

data "ibm_is_dedicated_hosts" "ds_dedicated_hosts" {
}

```

## Attribute Reference

The following attributes are exported:

* `dedicated_hosts` - List of all dedicated hosts in the IBM Cloud Infrastructure.
  * `id` - The id of the dedicated host.
  * `name` - The name of the dedicated host.
  * `group` - The dedicated host group name of this dedicated host.
  * `admin_state` - The administrative state of the dedicated host.
  * `available_memory` - The amount of memory in gibibytes that is currently available for instances.
  * `available_vcpu` - The available VCPU for the dedicated host.
  * `created_at` - The time (Created On) of the dedicated host.
  * `resource_group` - The resource group name of this dedicated host.
  * `lifecycle_state` - The lifecycle state of the dedicated host resource.
  * `crn` - The CRN for this dedicated host.
  * `href` - The URL for this dedicated host.
  * `instance_placement_enabled` - The instance placement enablement attribute. If set to true, instances can be placed on this dedicated host
  * `instances` - Instances that are allocated to this dedicated host.
  * `profile` - The profile name for this dedicated host profile.
  * `zone` - The name for this zone.
