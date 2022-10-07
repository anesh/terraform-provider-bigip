
# Overview

A [Terraform](terraform.io) provider for F5 BigIP LTM and BigIP GTM/DNS. This repo is a fork of the offical [repo](https://github.com/F5Networks/terraform-provider-bigip)
with the addition of a resource provider to configure BigIP GTM/DNS resources

# Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x / 0.12.x /0.13.x
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

# F5 BigIP LTM requirements

- This provider uses the iControlREST API, make sure that it is installed and enabled on your F5 device before proceeding.

These BIG-IP versions are supported in these Terraform versions.

| BIG-IP version	|Terraform 1.x  |	Terraform 0.13  |	Terraform 0.12  | Terraform 0.11  |
|-----------------|---------------|-----------------|-----------------|-----------------|
| BIG-IP 16.x	    |      X        |       X         |       X         |      X          |
| BIG-IP 15.x	    |      X        |       X         |       X         |      X          |
| BIG-IP 14.x	    | 	   X        |       X         |       X         |      X          |
| BIG-IP 12.x	    |      X        |      	X         |       X         |      X          | 
| BIG-IP 13.x	    |      X        |       X         |       X         |      X          |


# Documentation

Below is an example of how you can use the provider to create GTM/DNS resources

```
resource "bigip_gtm_datacenter" "datacenter" {
  name = "/Common/test"

}

resource "bigip_gtm_server" "server" {
  name = "/Common/Gw-cr-F5-2-lab.ctc"
  datacenter = "/Common/test"
  devices {
    name = "test4.ns.ctc"
    address = "5.5.5.5"

  }

}

resource "bigip_gtm_poola" "pool" {
  name = "/Common/test.ns.ctc"
  members = ["Gw-cr-F5-2-lab.ctc:test_vs"]
}

resource "bigip_gtm_wideipa" "wideip" {
  name = "/Common/test.wip.ns.ctc"
  pools {
    name = "test.ns.ctc"
    order = "0"
    ratio = "1"
  }
}

```


# Using the Provider

You can download the binary from the releases section of this repo and follow the instructions of installing terraform plugins.


