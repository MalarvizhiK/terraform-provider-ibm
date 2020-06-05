---
layout: "ibm"
page_title: "IBM : dedicated_host_groups"
sidebar_current: "docs-ibm-datasources-is-dedicated_host_groups"
description: |-
  List Dedicated host groups.
---

# ibm\_is_dedicated_host_groups

Import the details of an existing IBM dedicated host groups as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```hcl
resource "ibm_is_dedicated_host_group" "group1" {
  name           = "group1"
  resource_group = "42ee88169dd7406584925447bc5b25eb"
  zone = "us-south-1"
}


data "ibm_is_dedicated_host_groups" "ds_dedicated_host_groups" {
}

```

## Attribute Reference

The following attributes are exported:

* `groups` - List of all dedicated host groups in the IBM Cloud Infrastructure.
  * `id` - The id of the dedicated host group.
  * `name` - The name of the dedicated host group.
  * `created_at` - The time (Created On) of the dedicated host group.
  * `resource_group` - The resource group name of this dedicated host group.
  * `dedicated_hosts` - The dedicated hosts that are in this dedicated host group.
  * `crn` - The CRN for this dedicated host group.
  * `href` - The URL for this dedicated host group.
  * `zone` - The name for this zone.
