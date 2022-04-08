# output "kk"{
#     value= var.node_serial
# }
output "nn"{
    value = local.nodedata
}

output "ni"{
    value = local.nodeinfo
}
# output "mm"{
#     value = data.nutanix_foundation_central_imaged_node_details.nodedetails
# }

# output "cls"{
#     value = resource.nutanix_foundation_central_image_cluster.this
# }