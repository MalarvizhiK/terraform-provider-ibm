package ibm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isDedicatedHostGroupID         = "id"
	isDedicatedHostGroupCRN        = "crn"
	isDedicatedHostGroupHref       = "href"
	isDedicatedHostGrpName         = "name"
	isDedicatedHostGroupHosts      = "dedicated_hosts"
	isDedicatedHostResourceGrp     = "resource_group"
	isDedicatedHostGrpZone         = "zone"
	isDedicatedHostGroups          = "groups"
	isDedicatedHostGroupsCreatedAt = "created_at"
)

func dataSourceIBMISDedicatedHostGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISDedicatedHostGroupsRead,

		Schema: map[string]*schema.Schema{
			isDedicatedHostGroups: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated host groups",

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isDedicatedHostGrpName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of dedicatated host group name.",
						},
						isDedicatedHostGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host group ID.",
						},
						isDedicatedHostGroupCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host group CRN.",
						},
						isDedicatedHostGroupHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host group href.",
						},
						isDedicatedHostGroupsCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the dedicated host group was created.",
						},
						isDedicatedHostGroupHosts: {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						isDedicatedHostResourceGrp: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Host group resource group",
						},
						isDedicatedHostGrpZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Host group zone",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISDedicatedHostGroupsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listDedHostGrpsOptions := sess.NewListDedicatedHostGroupsOptions()

	result, resp, err := sess.ListDedicatedHostGroups(listDedHostGrpsOptions)
	if err != nil {
		log.Printf("Error reading list of dedicated host groups:%s", resp)
		return err
	}

	hostGrps := make([]map[string]interface{}, 0)

	for _, instance := range result.Groups {
		grp := map[string]interface{}{}
		grp[isDedicatedHostGrpName] = *instance.Name
		grp[isDedicatedHostGroupID] = *instance.ID
		grp[isDedicatedHostGroupCRN] = *instance.Crn
		grp[isDedicatedHostGroupHref] = *instance.Href
		grp[isDedicatedHostGroupsCreatedAt] = instance.CreatedAt.String()
		if len(instance.DedicatedHosts) != 0 {
			hostList := []string{}
			for i := 0; i < len(instance.DedicatedHosts); i++ {
				dedhost := instance.DedicatedHosts[i]
				hostList = append(hostList, string(*dedhost.ID))
			}
			grp[isDedicatedHostGroupHosts] = hostList
		}
		grp[isDedicatedHostResourceGrp] = *instance.ResourceGroup.ID
		grp[isDedicatedHostGrpZone] = *instance.Zone.Name
		hostGrps = append(hostGrps, grp)
	}

	d.SetId(dataSourceIBMISDedicatedHostGroupsID(d))
	d.Set(isDedicatedHostGroups, hostGrps)
	return nil
}

// dataSourceIBMISDedicatedHostGroupsID returns a reasonable ID for dedicated host groups.
func dataSourceIBMISDedicatedHostGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
