package ibm

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

func TestAccIBMISDedicatedHostGroupsDataSource_basic(t *testing.T) {
	log.Println("Before test case")
	node := "data.ibm_is_dedicated_host_groups.test1"
	hostgroupname := fmt.Sprintf("dedicatedhostgroup%d", acctest.RandIntRange(100, 200))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMISDedicatedHostGroupsDataSourceConfig(hostgroupname, ISZoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "groups.#"),
				),
			},
		},
	})
	log.Println("After test case")
}

func testAccCheckIBMISDedicatedHostGroupsDataSourceConfig(hostgroupname, zonename string) string {
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
	
    data "ibm_is_dedicated_host_groups" "test1" {
		# depends_on = [ibm_is_dedicated_host_group.group1]
	}`, hostgroupname, zonename)
}

func testAccCheckIBMISDedicatedHostGroupsDataSourceDestroy(s *terraform.State) error {
	instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_dedicated_host_group" {
			continue
		}
		log.Printf("Destroy called ...%s", rs.Primary.ID)

		delOptions := &vpcv1.DeleteDedicatedHostGroupOptions{
			ID: &rs.Primary.ID,
		}
		response, err := instanceC.DeleteDedicatedHostGroup(delOptions)

		if err != nil && response.StatusCode != 404 {
			log.Printf("Error deleting dedicated host group:%s", response)
			return err
		}

		getOptions := &vpcv1.GetDedicatedHostGroupOptions{
			ID: &rs.Primary.ID,
		}
		_, _, err = instanceC.GetDedicatedHostGroup(getOptions)

		if err == nil {
			return fmt.Errorf("instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
