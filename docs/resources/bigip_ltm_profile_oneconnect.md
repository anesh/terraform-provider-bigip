---
layout: "bigip"
page_title: "BIG-IP: bigip_ltm_profile_oneconnect"
subcategory: "Local Traffic Manager(LTM)"
description: |-
  Provides details about bigip_ltm_profile_oneconnect resource
---

# bigip\_ltm\_profile_oneconnect

`bigip_ltm_profile_oneconnect` Configures a custom profile_oneconnect for use by health checks.

Resources should be named with their "full path". The full path is the combination of the partition + name (example: /Common/my-pool ) or  partition + directory + name of the resource  (example: /Common/test/my-pool )

## Example Usage


```hcl
resource "bigip_ltm_profile_oneconnect" "test-oneconnect" {
  name = "/Common/test-oneconnect"
}

```      

## Argument Reference

* `name` (Required,`type string`) Name of Profile should be full path.The full path is the combination of the `partition + profile_name`,For example `/Common/test-oneconnect-profile`.

* `partition` - (Optional,`type string`) Displays the administrative partition within which this profile resides

* `defaults_from` - (Optional,`type string`) Specifies the profile that you want to use as the parent profile. Your new profile inherits all settings and values from the parent profile specified.

* `idle_timeout_override` - (Optional,`type string`) Specifies the number of seconds that a connection is idle before the connection flow is eligible for deletion. Possible values are `disabled`, `indefinite`, or a numeric value that you specify. The default value is `disabled`

* `limit_type` - (Optional,`type string`) Controls how connection limits are enforced in conjunction with OneConnect. The default is `None`. Supported Values: `[None,idle,strict]`

* `share_pools` - (Optional,`type string`) Specify if you want to share the pool, default value is `disabled`.

* `max_age` - (Optional,`type int`) Specifies the maximum age in number of seconds allowed for a connection in the connection reuse pool. For any connection with an age higher than this value, the system removes that connection from the reuse pool. The default value is `86400`.

* `max_reuse` - (Optional,`type int`) Specifies the maximum number of times that a server-side connection can be reused. The default value is `1000`.

* `max_size` - (Optional,`type int`) Specifies the maximum number of connections that the system holds in the connection reuse pool. If the pool is already full, then the server-side connection closes after the response is completed. The default value is `10000`.

* `source_mask` - (Optional,`type string`) Specifies a source IP mask. The default value is `0.0.0.0`. The system applies the value of this option to the source address to determine its eligibility for reuse. A mask of 0.0.0.0 causes the system to share reused connections across all clients. A host mask (all 1's in binary), causes the system to share only those reused connections originating from the same client IP address.


## Import

BIG-IP LTM oneconnect profiles can be imported using the `name` , e.g.

```
$ terraform import bigip_ltm_profile_oneconnect.test-oneconnect /Common/test-oneconnect
```
