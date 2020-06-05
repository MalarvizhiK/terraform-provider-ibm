---
layout: "ibm"
page_title: "IBM : dedicated_host_group"
sidebar_current: "docs-ibm-resource-is-dedicated_host_group"
description: |-
  Manages IBM IS dedicated_host_group.
---

# ibm\_is_dedicated_host_group

Provides a dedicated host group resource. This allows dedicated host group to be created, updated, and cancelled.


## Example Usage

```hcl

resource "ibm_is_dedicated_host_group" "test_dedicated_host_group" {
  name = "testdedicatedhostgroup"
  resource_group = "42ee88169dd7406584925447bc5b25eb"
  zone = "us-south-1"
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The dedicated host group name.
* `resource_group` - (Required, string) The resource group id of this dedicated host. 
* `zone` - (Required, string) The name of this zone. 

## Attribute Reference

The following attributes are exported:

* `id` - The id of the dedicated host.
* `name` - The name of the dedicated host.
* `created_at` - The time (Created On) of the dedicated host.
* `resource_group` - The resource group id of this dedicated host group.
* `crn` - The CRN for this dedicated host group.
* `href` - The URL for this dedicated host.
* `dedicated_hosts` - The dedicated hosts that are in this dedicated host group.
* `zone` - The name for this zone.

## Import

ibm_is_dedicated_host_group can be imported using dedicated host group ID, eg

```
$ terraform import ibm_is_dedicated_host_group.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
