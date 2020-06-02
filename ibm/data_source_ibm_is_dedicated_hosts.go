package ibm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dedicatedHosts                 = "dedicated_hosts"
	hostID                         = "id"
	instancePlacementEnabled       = "instance_placement_enabled"
	hostName                       = "name"
	hostCRN                        = "crn"
	hostHref                       = "href"
	adminState                     = "admin_state"
	availableMemory                = "available_memory"
	availableVCPU                  = "available_vcpu"
	lifeCycleState                 = "lifecycle_state"
	architecture                   = "architecture"
	count                          = "count"
	dedicatedHostGroup             = "group"
	dedicatedHostGroupID           = "id"
	dedicatedHostGroupCRN          = "crn"
	dedicatedHostGroupHref         = "href"
	dedicatedHostGroupName         = "name"
	dedicatedHostInstances         = "instances"
	dedicatedHostInstanceID        = "id"
	dedicatedHostInstanceCRN       = "crn"
	dedicatedHostInstanceHref      = "href"
	dedicatedHostInstanceName      = "name"
	dedicatedHostProfile           = "profile"
	dedicatedHostProfileHref       = "href"
	dedicatedHostProfileName       = "name"
	dedicatedHostResourceGroup     = "resource_group"
	dedicatedHostResourceGroupID   = "id"
	dedicatedHostResourceGroupHref = "href"
	dedicatedHostResourceGroupName = "name"
	dedicatedHostZone              = "zone"
	dedicatedHostZoneHref          = "href"
	dedicatedHostZoneName          = "name"

	lifecycleStatePending  = "pending"
	lifecycleStateStable   = "stable"
	lifecycleStateUpdating = "updating"
	lifecycleStateWaiting  = "waiting"
	lifecycleStateFailed   = "failed"
	lifecycleStateDeleting = "deleting"
)

func dataSourceIBMisDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMisDedicatedHostRead,

		Schema: map[string]*schema.Schema{
			dedicatedHosts: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated hosts",

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						hostName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of dedicatated host.",
						},
						instancePlacementEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, instances can be placed on this dedicated host.",
						},
						hostID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host ID.",
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
							Computed:    true,
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
							Computed:    true,
							Description: "The dedicated host profile.",
						},
						dedicatedHostResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host resource group.",
						},
						dedicatedHostZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMisDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listDedicatedHostsOptions := sess.NewListDedicatedHostsOptions()

	result, resp, err := sess.ListDedicatedHosts(listDedicatedHostsOptions)
	if err != nil {
		log.Printf("Error reading list of dedicated hosts:%s", resp)
		return err
	}

	hosts := make([]map[string]interface{}, 0)

	for _, instance := range result.DedicatedHosts {
		instances := []string{}
		for _, ins := range instance.Instances {
			instances = append(instances, *ins.ID)
			log.Println("host instance id", *ins.ID)
		}
		host := map[string]interface{}{
			hostName:                   instance.Name,
			hostID:                     instance.ID,
			hostCRN:                    instance.Crn,
			hostHref:                   instance.Href,
			adminState:                 instance.AdminState,
			availableMemory:            instance.AvailableMemory,
			instancePlacementEnabled:   instance.InstancePlacementEnabled,
			lifeCycleState:             instance.LifecycleState,
			availableVCPU:              instance.AvailableVcpu.Count,
			dedicatedHostGroup:         instance.Group.Name,
			dedicatedHostResourceGroup: instance.ResourceGroup.Name,
			dedicatedHostZone:          instance.Zone.Name,
			dedicatedHostProfile:       instance.Profile.Name,
			dedicatedHostInstances:     instances,
		}

		hosts = append(hosts, host)
	}

	d.SetId(dataSourceIBMisDedicatedHostID(d))
	log.Println("hosts list ", hosts)
	d.Set(dedicatedHosts, hosts)
	return nil
}

// dataSourceIBMisDedicatedHostID returns a reasonable ID for dedicated hosts.
func dataSourceIBMisDedicatedHostID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
