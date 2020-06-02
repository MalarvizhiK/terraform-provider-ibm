package ibm

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

func TestAccIBMISInstanceTemplatesDataSource_basic(t *testing.T) {
	var instance string
	node := "data.ibm_is_instance_templates.test1"
	vpcname := fmt.Sprintf("vpc-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("name-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("vpc-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("vpc-ssh-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceTemplatesDataSourceConfig(vpcname, subnetname, sshname, publicKey, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISInstanceTemplateExists("ibm_is_instance_template.testacc_instance", instance),
					resource.TestCheckResourceAttrSet(node, "templates.0.name"),
					resource.TestCheckResourceAttrSet(node, "templates.0.id"),
					resource.TestCheckResourceAttrSet(node, "templates.0.crn"),
				),
			},
		},
	})
}

func testAccCheckIBMISInstanceTemplatesDataSourceConfig(vpcname, subnetname, sshname, publicKey, name string) string {
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		name = "Default"
	}

	resource "ibm_is_vpc" "testacc_vpc" {
		depends_on = [data.ibm_resource_group.rg]
		name = "%s"
		resource_group = data.ibm_resource_group.rg.id
	  }
	  
	  resource "ibm_is_subnet" "testacc_subnet" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
	  }
	  
	  resource "ibm_is_ssh_key" "testacc_sshkey" {
		name       = "%s"
		public_key = "%s"
	  }
	  
	  resource "ibm_is_instance_template" "testacc_instance" {
		name    = "%s"
		image   = "%s"
		profile = "%s"
		primary_network_interface {
		  port_speed = "100"
		  subnet     = ibm_is_subnet.testacc_subnet.id
		}
		vpc  = ibm_is_vpc.testacc_vpc.id
		zone = "%s"
		keys = [ibm_is_ssh_key.testacc_sshkey.id]
		network_interfaces {
		  subnet = ibm_is_subnet.testacc_subnet.id
		  name   = "eth1"
		}
	  }
	 
	  data "ibm_is_instance_templates" "test1" {
		depends_on = [ibm_is_instance_template.testacc_instance]
	}
	  
	  `, vpcname, subnetname, ISZoneName, ISCIDR, sshname, publicKey, name, isImage, instanceProfileName, ISZoneName)
}

func testAccCheckIBMISInstanceTemplateExists(n string, instance string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}
		instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
		getinsOptions := &vpcv1.GetInstanceTemplateOptions{
			ID: &rs.Primary.ID,
		}
		instanceIntf, _, err := instanceC.GetInstanceTemplate(getinsOptions)
		instance1 := instanceIntf.(*vpcv1.InstanceTemplate)
		if err != nil {
			return err
		}
		instance = *instance1.ID

		return nil
	}
}

func testAccCheckIBMISInstanceTemplateDestroy(s *terraform.State) error {
	instanceC, _ := testAccProvider.Meta().(ClientSession).VpcV1API()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_instance_template" {

			log.Printf("Destroy called ...%s", rs.Primary.ID)

			deleteinstanceTemplateOptions := &vpcv1.DeleteInstanceTemplateOptions{
				ID: &rs.Primary.ID,
			}
			response, err := instanceC.DeleteInstanceTemplate(deleteinstanceTemplateOptions)

			if err != nil && response.StatusCode != 404 {
				log.Printf("Error deleting instance template:%s", response)
				return err
			}

			getinsOptions := &vpcv1.GetInstanceTemplateOptions{
				ID: &rs.Primary.ID,
			}
			_, response, err = instanceC.GetInstanceTemplate(getinsOptions)

			if err == nil {
				return fmt.Errorf("instance still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}
