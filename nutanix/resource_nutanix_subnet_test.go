package nutanix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnet_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 500)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.next-iac-managed"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_subnet" {
			continue
		}
		if _, err := resourceNutanixSubnetExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccNutanixSubnetConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
		
resource "nutanix_subnet" "next-iac-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information
  name        = "next-iac-managed-%d"
  vlan_id     = %d
  subnet_type = "VLAN"

  # Managed L3 Networks
  # This bit is only needed if you intend to turn on IPAM
  prefix_length = 20

  default_gateway_ip = "10.5.80.1"
  subnet_ip          = "10.5.80.0"

  #dhcp_options {
  #    boot_file_name   = "bootfile"
  #    tftp_server_name = "1.2.3.200"
  #    domain_name      = "nutanix"
  #}

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}
`, r, r)
}