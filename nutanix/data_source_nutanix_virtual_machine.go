package nutanix

import (
	"fmt"
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixVirtualMachineRead,

		Schema: getDataSourceVMSchema(),
	}
}

func dataSourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	vm, ok := d.GetOk("vm_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetVM(vm.(string))
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)
	n := setNicList(resp.Status.Resources.NicList)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", getReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
		return err
	}
	if err := d.Set("nic_list", n); err != nil {
		return err
	}
	if err := d.Set("host_reference", getReferenceValues(resp.Status.Resources.HostReference)); err != nil {
		return err
	}
	if err := d.Set("nutanix_guest_tools", setNutanixGuestTools(resp.Status.Resources.GuestTools)); err != nil {
		return err
	}
	if err := d.Set("gpu_list", setGPUList(resp.Status.Resources.GpuList)); err != nil {
		return err
	}
	if err := d.Set("parent_reference", getReferenceValues(resp.Status.Resources.ParentReference)); err != nil {
		return err
	}

	diskAddress := make(map[string]interface{})
	mac := ""
	b := make([]string, 0)

	if resp.Status.Resources.BootConfig != nil {
		if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
			i := strconv.Itoa(int(utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))
			diskAddress["device_index"] = i
			diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		}
		if resp.Status.Resources.BootConfig.BootDeviceOrderList != nil {
			b = utils.StringValueSlice(resp.Status.Resources.BootConfig.BootDeviceOrderList)
		}
		mac = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
	}

	d.Set("boot_device_order_list", b)
	d.Set("boot_device_disk_address", diskAddress)
	d.Set("boot_device_mac_address", mac)

	sysprep := make(map[string]interface{})
	sysrepCV := make(map[string]string)
	cloudInit := make(map[string]interface{})
	cloudInitCV := make(map[string]string)
	isOv := false
	if resp.Status.Resources.GuestCustomization != nil {
		isOv = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)
		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			if resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues {
					cloudInitCV[k] = v
				}
			}
		}
		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)

			if resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues {
					sysrepCV[k] = v
				}
			}
		}
	}
	if err := d.Set("guest_customization_cloud_init_custom_key_values", cloudInitCV); err != nil {
		return err
	}
	if err := d.Set("guest_customization_sysprep_custom_key_values", sysrepCV); err != nil {
		return err
	}
	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return err
	}
	if err := d.Set("guest_customization_cloud_init", cloudInit); err != nil {
		return err
	}

	d.Set("hardware_clock_timezone", utils.StringValue(resp.Status.Resources.HardwareClockTimezone))
	d.Set("cluster_reference_name", utils.StringValue(resp.Status.ClusterReference.Name))
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes))
	d.Set("guest_os_id", utils.StringValue(resp.Status.Resources.GuestOsID))
	d.Set("power_state", utils.StringValue(resp.Status.Resources.PowerState))
	d.Set("num_vcpus_per_socket", utils.Int64Value(resp.Status.Resources.NumVcpusPerSocket))
	d.Set("num_sockets", utils.Int64Value(resp.Status.Resources.NumSockets))
	d.Set("memory_size_mib", utils.Int64Value(resp.Status.Resources.MemorySizeMib))
	d.Set("guest_customization_is_overridable", isOv)
	d.Set("should_fail_on_script_failure", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure))
	d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec))
	d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism))
	d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled))
	d.SetId(*resp.Metadata.UUID)

	return d.Set("disk_list", setDiskList(resp.Status.Resources.DiskList))
}

func getDataSourceVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vm_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// COMPUTED
		"message_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": {
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hypervisor_type": {
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"nic_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"floating_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_function_nic_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ip_endpoint_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"subnet_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"guest_os_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"power_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nutanix_guest_tools": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"available_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"iso_mount_state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guest_os_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled_capability_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vss_snapshot_capable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"is_reachable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vm_mobility_drivers_installed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"num_vcpus_per_socket": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"num_sockets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gpu_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frame_buffer_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vendor": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"fraction": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mode": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"num_virtual_display_heads": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guest_driver_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"parent_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"memory_size_mib": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"boot_device_order_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device_disk_address": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_index": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"adapter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"boot_device_mac_address": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"hardware_clock_timezone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guest_customization_cloud_init": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"meta_data": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"user_data": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_cloud_init_custom_key_values": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"guest_customization_is_overridable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guest_customization_sysprep": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"install_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"unattend_xml": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_sysprep_custom_key_values": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"should_fail_on_script_failure": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"enable_script_exec": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"power_state_mechanism": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vga_console_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disk_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"disk_size_bytes": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"device_properties": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"disk_address": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_index": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"adapter_type": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"data_source_reference": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"volume_group_reference": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}