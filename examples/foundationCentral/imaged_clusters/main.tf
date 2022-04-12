// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta.2
terraform {
    required_providers {
      nutanix = {
          source = "nutanix/nutanix"
          version = ">1.5.0-beta.2"
      }
    }
}

provider "nutanix" {
    username  = "user"
    password  = "pass"
    endpoint  = "10.x.xx.xx"
    insecure  = true
    port      = 9440
}


// resource to image and create cluster

resource "nutanix_foundation_central_image_cluster" "res" {
  cluster_name = "test-fc"
  common_network_settings {
    cvm_dns_servers = [
      "10.4.8.15"
    ]
    hypervisor_dns_servers = [
      "10.4.8.15"
    ]
    cvm_ntp_servers = [
      "0.pool.ntp.org"
    ]
    hypervisor_ntp_servers = [
      "0.pool.ntp.org"
    ]
  }
  redundancy_factor = 2
  node_list {
    cvm_gateway                   = "10.2.240.1"
    cvm_netmask                   = "255.255.240.0"
    cvm_ip                        = "10.2.243.76"
    hypervisor_gateway            = "10.2.240.1"
    hypervisor_netmask            = "255.255.240.0"
    hypervisor_ip                 = "10.2.243.72"
    hypervisor_hostname           = "IBIS19-3"
    imaged_node_uuid              = "e7c53f45-2e49-4784-57ea-c350cdd18fd0"
    use_existing_network_settings = false
    ipmi_gateway                  = "10.2.128.1"
    ipmi_netmask                  = "255.255.240.0"
    ipmi_ip                       = "10.2.133.217"
    image_now                     = true
    hypervisor_type               = "kvm"
    hardware_attributes_override = {
      default_workload      = "vdi"
      lcm_family            = "smc_gen_10"
      maybe_1GbE_only       = true
      robo_mixed_hypervisor = true
    }

  }
  node_list {
    cvm_gateway                   = "10.2.240.1"
    cvm_netmask                   = "255.255.240.0"
    cvm_ip                        = "10.2.243.74"
    hypervisor_gateway            = "10.2.240.1"
    hypervisor_netmask            = "255.255.240.0"
    hypervisor_ip                 = "10.2.243.70"
    hypervisor_hostname           = "IBIS19-1"
    imaged_node_uuid              = "c1542802-bfe4-4878-6f03-dfaf0f4b7ade"
    use_existing_network_settings = false
    ipmi_gateway                  = "10.2.128.1"
    ipmi_netmask                  = "255.255.240.0"
    ipmi_ip                       = "10.2.133.215"
    image_now                     = true
    hypervisor_type               = "kvm"
  }
  node_list {
    cvm_gateway                   = "10.2.240.1"
    cvm_netmask                   = "255.255.240.0"
    cvm_ip                        = "10.2.243.75"
    hypervisor_gateway            = "10.2.240.1"
    hypervisor_netmask            = "255.255.240.0"
    hypervisor_ip                 = "10.2.243.71"
    hypervisor_hostname           = "IBIS19-2"
    imaged_node_uuid              = "7e7f071a-2673-4e65-5b2e-a026ba0873d4"
    use_existing_network_settings = false
    ipmi_gateway                  = "10.2.128.1"
    ipmi_netmask                  = "255.255.240.0"
    ipmi_ip                       = "10.2.133.216"
    image_now                     = true
    hypervisor_type               = "kvm"
  }
  aos_package_url = "http://endor.dyn.nutanix.com/builds/nos-builds/fraser-6.1.1-stable/08e579e4287942f28c91436c30712901eacfbeab/x86_64/release/tar/nutanix_installer_package-release-fraser-6.1.1-stable-08e579e4287942f28c91436c30712901eacfbeab-x86_64.tar.gz"
}

output "res1"{
    value = resource.nutanix_foundation_central_image_cluster.res
}


// datasource to List all the clusters created using Foundation Central.

// data "nutanix_foundation_central_imaged_clusters_list" "cls" {}

// output "cls1" {
//   value = data.nutanix_foundation_central_imaged_clusters_list.cls
// }

// // datasource to Get a cluster created using Foundation Central.
// data "nutanix_foundation_central_cluster_details" "clsDet" {
//   imaged_cluster_uuid = "<imaged_cluster_uuid>"
// }

// output "clsDet1" {
//   value = data.nutanix_foundation_central_cluster_details.clsDet
// }