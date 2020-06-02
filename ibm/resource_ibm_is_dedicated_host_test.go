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

func TestAccIBMISDedicatedHost_basic(t *testing.T) {
	log.Println("Before test case")
	hostgroupname := fmt.Sprintf("dedicatedhostgroup%d", acctest.RandIntRange(100, 200))
	hostname := fmt.Sprintf("dedicatedhost%d", acctest.RandIntRange(100, 200))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMISDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISDedicatedHostConfig(hostgroupname, hostname, ISZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISDedicatedHostExists("ibm_is_dedicated_host.host1", hostname),
					resource.TestCheckResourceAttr(
						"ibm_is_dedicated_host.host1", "name", hostname),
					resource.TestCheckResourceAttr(
						"ibm_is_dedicated_host.host1", "group", hostgroupname),
				),
			},
		},
	})
	log.Println("After test case")
}

func testAccCheckIBMISDedicatedHostDestroy(s *terraform.State) error {

	instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ibm_is_dedicated_host" {

			getOptions := &vpcv1.GetDedicatedHostOptions{
				ID: &rs.Primary.ID,
			}
			_, _, err := instanceC.GetDedicatedHost(getOptions)

			if err == nil {
				return fmt.Errorf("instance still exists: %s", rs.Primary.ID)
			}
		}
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ibm_is_dedicated_host_group" {

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

func testAccCheckIBMISDedicatedHostExists(n string, instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
		getOptions := &vpcv1.GetDedicatedHostOptions{
			ID: &rs.Primary.ID,
		}
		foundins, _, err := instanceC.GetDedicatedHost(getOptions)
		if err != nil {
			return err
		}
		instance = *foundins.ID
		return nil
	}
}

func testAccCheckIBMISDedicatedHostConfig(hostgroupname, hostname, zonename string) string {
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
		instance_placement_enabled = "true"
		profile = "dh2-56x464"
	} 

   `, hostgroupname, zonename, hostname)
}
