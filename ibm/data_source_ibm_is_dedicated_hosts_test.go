package ibm

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMISDedicatedHostsDataSource_basic(t *testing.T) {
	log.Println("Before test case")
	log.Printf("Zone name %s", ISZoneName)
	node := "data.ibm_is_dedicated_hosts.test1"
	hostgroupname := fmt.Sprintf("dedicatedhostgroup%d", acctest.RandIntRange(100, 200))
	hostname := fmt.Sprintf("dedicatedhost%d", acctest.RandIntRange(100, 200))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMISDedicatedHostsDataSourceConfig(hostgroupname, hostname, ISZoneName),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(node, "groups"),
					resource.TestCheckResourceAttrSet(node, "dedicated_hosts.#"),
				),
			},
		},
	})
	log.Println("After test case")
}

func testAccCheckIBMISDedicatedHostsDataSourceConfig(hostgroupname, hostname, zonename string) string {
	log.Println("Inside test case")
	// status filter defaults to empty
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		name = "Default"
	}	

	resource "ibm_is_dedicated_host_group" "group1" {
		name           = "%s"
		resource_group = data.ibm_resource_group.rg.id
		zone = "%s"
	}
	
	resource "ibm_is_dedicated_host" "host1" {
		name           = "%s"
		resource_group = data.ibm_resource_group.rg.id
		group = ibm_is_dedicated_host_group.group1.id
		instance_placement_enabled = "false"
		profile = "dh2-56x464"
	}
	
	data "ibm_is_dedicated_hosts" "test1" {
		
	}`, hostgroupname, zonename, hostname)
}
