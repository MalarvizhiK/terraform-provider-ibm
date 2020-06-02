package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

var dedicatedHostV2PackageType = "IS_DEDICATED_HOST"

func resourceIBMisDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisDedicatedHostCreate,
		Read:     resourceIBMisDedicatedHostRead,
		Delete:   resourceIBMisDedicatedHostDelete,
		Exists:   resourceIBMisDedicatedHostExists,
		Update:   resourceIBMisDedicatedHostUpdate,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			hostName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of dedicatated host.",
			},
			instancePlacementEnabled: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "If set to true, instances can be placed on this dedicated host.",
			},
			hostCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The dedicated host CRN.",
			},
			hostHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The dedicated host href.",
			},
			adminState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the dedicated host.",
			},
			availableMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory in gibibytes that is currently available for instances.",
			},
			lifeCycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dedicated host resource.",
			},
			availableVCPU: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dedicated host group this dedicated host is in.",
			},
			dedicatedHostGroup: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host group.",
			},
			dedicatedHostInstances: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The dedicated host instances.",
			},
			dedicatedHostProfile: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host profile.",
			},
			dedicatedHostResourceGroup: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host resource group.",
			},
		},
	}
}

func dedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	grpID := d.Get(dedicatedHostGroup).(string)

	grpTemplate := &vpcv1.DedicatedHostPrototypeGroup{
		ID: &grpID,
	}
	profile := d.Get(dedicatedHostProfile).(string)
	profileTemplate := &vpcv1.DedicatedHostPrototypeProfile{
		Name: &profile,
	}
	name := d.Get(hostName).(string)
	instancePlacementEnabled := d.Get(instancePlacementEnabled).(bool)
	resourceGrp := d.Get(dedicatedHostResourceGroup).(string)
	resourceGrpTemplate := &vpcv1.ResourceGroupIdentity{
		ID: &resourceGrp,
	}
	options := &vpcv1.CreateDedicatedHostOptions{
		Group:                    grpTemplate,
		Profile:                  profileTemplate,
		Name:                     &name,
		InstancePlacementEnabled: &instancePlacementEnabled,
		ResourceGroup:            resourceGrpTemplate,
	}

	result, resp, err := sess.CreateDedicatedHost(options)
	if err != nil {
		log.Printf("[DEBUG] Dedicated host error %s\n%s", err, resp)
		return err
	}
	d.SetId(*result.ID)

	log.Printf("[INFO] dedicate host : %s", *result.ID)
	return nil
}

func dedicatedHostGet(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()
	options := sess.NewGetDedicatedHostOptions(id)
	result, resp, err := sess.GetDedicatedHost(options)
	if err != nil {
		log.Printf("Error getting dedicated hosts:%s", resp)
		return err
	}
	d.Set(hostName, result.Name)
	d.Set(hostID, result.ID)
	d.Set(instancePlacementEnabled, result.InstancePlacementEnabled)
	d.Set(hostCRN, result.Crn)
	d.Set(hostHref, result.Href)
	d.Set(adminState, result.AdminState)
	d.Set(availableMemory, result.AvailableMemory)
	d.Set(lifeCycleState, result.LifecycleState)
	d.Set(dedicatedHostGroup, result.Group.Name)
	d.Set(dedicatedHostProfile, result.Profile.Name)
	d.Set(dedicatedHostZone, result.Zone.Name)
	return nil
}

func dedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()
	name := d.Get(hostName).(string)
	instancePlacementEnabled := d.Get(instancePlacementEnabled).(bool)

	options := sess.NewUpdateDedicatedHostOptions(id)
	options.SetName(name)
	options.SetInstancePlacementEnabled(instancePlacementEnabled)

	result, resp, err := sess.UpdateDedicatedHost(options)
	if err != nil {
		log.Printf("[DEBUG] Update Dedicated host error %s\n%s", err, resp)
		return err
	}
	d.SetId(*result.ID)

	log.Printf("[INFO] dedicate host update : %s", *result.ID)
	return nil
}

func dedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()

	options := sess.NewDeleteDedicatedHostOptions(id)

	resp, err := sess.DeleteDedicatedHost(options)
	if err != nil {
		log.Printf("[DEBUG] Delete Dedicated host error %s\n%s", err, resp)
		return err
	}

	log.Printf("[INFO] dedicate host deleted")
	return nil
}

func isWaitForDedicatedHostAvailable(instanceC *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for instance (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", lifecycleStateWaiting},
		Target:     []string{lifecycleStateStable, lifecycleStateFailed, "", ""},
		Refresh:    isInstanceRefreshFunc(instanceC, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isDedicatedHostRefreshFunc(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		options := instanceC.NewGetDedicatedHostOptions(id)
		instance, response, err := instanceC.GetDedicatedHost(options)
		if err != nil {
			return nil, "", fmt.Errorf("Error getting dedicated host instance: %s\n%s", err, response)
		}
		d.Set(lifeCycleState, *instance.LifecycleState)

		if *instance.LifecycleState == lifecycleStateStable || *instance.LifecycleState == lifecycleStateFailed || *instance.LifecycleState == lifecycleStateUpdating {
			return instance, *instance.LifecycleState, nil
		}

		return instance, lifecycleStateUpdating, nil
	}
}

func resourceIBMisDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	err := dedicatedHostCreate(d, meta)
	if err != nil {
		return err
	}
	return resourceIBMisDedicatedHostRead(d, meta)
}

func resourceIBMisDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {

	err := dedicatedHostGet(d, meta)
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	return nil
}

func resourceIBMisDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {

	err := dedicatedHostUpdate(d, meta)
	if err != nil {
		return fmt.Errorf("Error: dedicated host update failed %s", err)
	}
	return resourceIBMisDedicatedHostRead(d, meta)
}

func resourceIBMisDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	d.Set(instancePlacementEnabled, false)
	err := dedicatedHostUpdate(d, meta)
	if err != nil {
		return fmt.Errorf("Error: dedicated host update failed %s", err)
	}
	err = dedicatedHostDelete(d, meta)
	if err != nil {
		return fmt.Errorf("Error: dedicated host delete failed %s", err)
	}
	return nil
}

func resourceIBMisDedicatedHostExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	id := d.Id()
	options := sess.NewGetDedicatedHostOptions(id)
	result, resp, err := sess.GetDedicatedHost(options)
	if err != nil {
		log.Printf("Error getting dedicated hosts %s", resp)
		return false, err
	}
	return result.ID != nil && *result.ID == id, nil
}
