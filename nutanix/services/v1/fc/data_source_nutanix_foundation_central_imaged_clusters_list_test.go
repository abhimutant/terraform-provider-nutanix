package foundationcentral_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCClusterListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCClusterListDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_clusters_list.cls", "imaged_clusters.#"),
				),
			},
		},
	})
}

func testAccFCClusterListDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_clusters_list" "cls" {}
	`
}