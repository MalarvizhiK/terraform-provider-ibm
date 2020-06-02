package ibm

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

func TestAccIBMISDedicatedHostGroup_basic(t *testing.T) {
	log.Println("Before test case")
	hostgroupname := fmt.Sprintf("dedicatedhostgroup%d", acctest.RandIntRange(100, 200))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMISDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISDedicatedHostGroupConfig(hostgroupname, ISZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISDedicatedHostGroupExists("ibm_is_dedicated_host_group.group1", hostgroupname),
					resource.TestCheckResourceAttr(
						"ibm_is_dedicated_host_group.group1", "name", hostgroupname),
					resource.TestCheckResourceAttr(
						"ibm_is_dedicated_host_group.group1", "zone", ISZoneName),
				),
			},
		},
	})
	log.Println("After test case")
}

func testAccCheckIBMISDedicatedHostGroupDestroy(s *terraform.State) error {

	instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_dedicated_host_group" {
			getOptions := &vpcv1.GetDedicatedHostGroupOptions{
				ID: &rs.Primary.ID,
			}
			_, _, err := instanceC.GetDedicatedHostGroup(getOptions)

			if err == nil {
				return fmt.Errorf("instance still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIBMISDedicatedHostGroupExists(n string, instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
		getOptions := &vpcv1.GetDedicatedHostGroupOptions{
			ID: &rs.Primary.ID,
		}
		foundins, _, err := instanceC.GetDedicatedHostGroup(getOptions)
		if err != nil {
			return err
		}
		instance = *foundins.ID
		return nil
	}
}

func testAccCheckIBMISDedicatedHostGroupConfig(hostgroupname, zonename string) string {
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
	
   `, hostgroupname, zonename)
}
