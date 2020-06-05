---
layout: "ibm"
page_title: "IBM : dedicated_host"
sidebar_current: "docs-ibm-resource-is-dedicated_host"
description: |-
  Manages IBM IS Dedicated_host.
---

# ibm\_is_dedicated_host

Provides a dedicated host resource. This allows instance to be created, updated, and cancelled.


## Example Usage

```hcl

resource "ibm_is_dedicated_host" "test_dedicated_host" {
  name = "testdedicatedhost"
  instance_placement_enabled = "false"
  group = ibm_is_dedicated_host_group.test_dedicated_host_group.id
  resource_group = data.ibm_resource_group.rg.id
  profile = "dh2-56x464"
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The dedicated group name.
* `instance_placement_enabled` - (Required, string) The instance placement enablement attribute. If set to true, instances can be placed on this dedicated host. 
* `group` - (Required, string) The dedicated host group id of this dedicated host.
* `resource_group` - (Required, string) The resource group id of this dedicated host. 
* `profile` - (Required, string) The dedicated profile name. 

## Attribute Reference

The following attributes are exported:

* `id` - The id of the dedicated host.
* `name` - The name of the dedicated host.
* `group` - The dedicated host group id of this dedicated host.
* `admin_state` - The administrative state of the dedicated host.
* `available_memory` - The amount of memory in gibibytes that is currently available for instances.
* `available_vcpu` - The available VCPU for the dedicated host.
* `created_at` - The time (Created On) of the dedicated host.
* `resource_group` - The resource group id of this dedicated host.
* `lifecycle_state` - The lifecycle state of the dedicated host resource.
* `crn` - The CRN for this dedicated host.
* `href` - The URL for this dedicated host.
* `instance_placement_enabled` - The instance placement enablement attribute. If set to true, instances can be placed on this dedicated host
* `instances` - Instances that are allocated to this dedicated host.
* `profile` - The profile name for this dedicated host profile. 
* `zone` - The name for this zone.

## Import

ibm_is_dedicated_host can be imported using dedicated host ID, eg

```
$ terraform import ibm_is_dedicated_host.example d7bec597-4726-451f-8a63-e62e6f19c32c
```
