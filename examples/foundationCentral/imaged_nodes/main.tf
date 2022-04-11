// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta
terraform{
    required_providers{
        nutanix = {
            source = "nutanixtemp/nutanix"
            version = "1.99.99"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    username  = "admin"
    password  = "Nutanix.123"
    endpoint  = "10.2.242.13"
    insecure  = true
    port      = 9440
}

// datasource to List all the nodes registered with Foundation Central.
data "nutanix_foundation_central_imaged_nodes_list" "img"{}

output "img1"{
    value = data.nutanix_foundation_central_imaged_nodes_list.img
}

// datasource to Get the details of a single node given its UUID.
data "nutanix_foundation_central_imaged_node_details" "imgdet"{
    imaged_node_uuid = "<imaged_node_uuid>"

output "imgdetails"{
    value = data.nutanix_foundation_central_imaged_node_details.imgdet
}