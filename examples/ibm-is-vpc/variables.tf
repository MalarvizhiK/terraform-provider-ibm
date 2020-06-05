variable "zone1" {
  default = "us-south-1"
}

variable "zone2" {
  default = "us-south-2"
}

variable "ssh_public_key" {
  default = "~/.ssh/id_rsa.pub"
}

variable "image" {
  default = "7eb4e35b-4257-56f8-d7da-326d85452591"
}

variable "profile" {
  default = "bc1-2x8"
}

variable "region" {
  default     = "us-south"
  description = "The VPC Region that you want your VPC, networks and the virtual server to be provisioned in. To list available regions, run `ibmcloud is regions`."
}

variable "resource_group" {
  default     = "Default"
  description = "The resource group to use. If unspecified, the account's default resource group is used."
}

