package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

const (
	isDedicatedHostGroupName     = "name"
	isDedicatedHostGroupZone     = "zone"
	isDedicatedHostResourceGroup = "resource_group"
)

func resourceIBMISDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisDedicatedHostGroupCreate,
		Read:     resourceIBMisDedicatedHostGroupRead,
		Update:   resourceIBMisDedicatedHostGroupUpdate,
		Delete:   resourceIBMisDedicatedHostGroupDelete,
		Exists:   resourceIBMisDedicatedHostGroupExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isDedicatedHostGroupName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateISName,
				Description:  "Dedicated Host Group name",
			},

			isDedicatedHostGroupZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Zone name",
			},

			isDedicatedHostResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Dedicated Host resource group",
			},
		},
	}
}

func resourceIBMisDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	name := d.Get(isDedicatedHostGroupName).(string)
	group := d.Get(isDedicatedHostResourceGroup).(string)
	zone := d.Get(isDedicatedHostGroupZone).(string)

	options := &vpcv1.CreateDedicatedHostGroupOptions{
		Name: &name,
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &group,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
	}

	hostGrp, response, err := sess.CreateDedicatedHostGroup(options)
	if err != nil {
		log.Printf("[DEBUG] Create Dedicated Host Group err %s\n%s", err, response)
		return err
	}
	d.SetId(*hostGrp.ID)

	log.Printf("[INFO] Dedicated Host Group Template : %s", *hostGrp.ID)

	if err != nil {
		log.Printf("Error creating dedicated host group:%s", response)
		return err
	}

	return resourceIBMisDedicatedHostGroupUpdate(d, meta)
}

func resourceIBMisDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {

	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()

	getOptions := &vpcv1.GetDedicatedHostGroupOptions{
		ID: &ID,
	}
	hostGrp, response, err := instanceC.GetDedicatedHostGroup(getOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Dedicated Host Group: %s\n%s", err, response)
	}
	d.Set(isDedicatedHostGroupName, *hostGrp.Name)

	if hostGrp.ResourceGroup != nil {
		d.Set(isDedicatedHostResourceGroup, *hostGrp.ResourceGroup.ID)
	}

	if hostGrp.Zone != nil {
		d.Set(isDedicatedHostGroupZone, *hostGrp.Zone.Name)
	}

	return nil
}

func resourceIBMisDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()
	getOptions := &vpcv1.GetDedicatedHostGroupOptions{
		ID: &ID,
	}
	_, detail, err := instanceC.GetDedicatedHostGroup(getOptions)
	if err != nil {
		log.Printf("Error fetching dedicated host group:%s", detail)
		return err
	}

	if d.HasChange(isDedicatedHostGroupName) {
		name := d.Get(isDedicatedHostGroupName).(string)
		updoptions := &vpcv1.UpdateDedicatedHostGroupOptions{
			ID:   &ID,
			Name: &name,
		}
		_, detail, err = instanceC.UpdateDedicatedHostGroup(updoptions)
		if err != nil {
			log.Printf("Error updating dedicated host group:%s", detail)
			return err
		}
	}

	return resourceIBMisDedicatedHostGroupRead(d, meta)
}

func resourceIBMisDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {

	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()
	delOptions := &vpcv1.DeleteDedicatedHostGroupOptions{
		ID: &ID,
	}
	response, err := instanceC.DeleteDedicatedHostGroup(delOptions)

	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting dedicated host group:%s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMisDedicatedHostGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	ID := d.Id()

	getOptions := &vpcv1.GetDedicatedHostGroupOptions{
		ID: &ID,
	}
	_, response, err := instanceC.GetDedicatedHostGroup(getOptions)
	if err != nil && response.StatusCode != 404 {
		return false, fmt.Errorf("Error Getting dedicated host group : %s\n%s", err, response)
	}
	if response.StatusCode == 404 {
		d.SetId("")
		return false, nil
	}
	return true, nil
}
