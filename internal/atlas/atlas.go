package atlas

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/atlas-sdk/v20231115014/admin"
)

func NewClient(apiKey, apiSecret string) (*admin.APIClient, error) {
	sdk, err := admin.NewClient(admin.UseDigestAuth(apiKey, apiSecret))
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func GetProjectByName(ctx context.Context, sdk *admin.APIClient, projectName string) (*Project, error) {
	project, _, err := sdk.ProjectsApi.GetProjectByName(ctx, projectName).Execute()
	if err != nil {
		return nil, err
	}

	return &Project{
		Id: *project.Id,
	}, nil
}

func GetCluster(ctx context.Context, projectId Project, clusterName string, sdk *admin.APIClient) (*admin.AdvancedClusterDescription, error) {
	cluster, _, err := sdk.ClustersApi.GetCluster(ctx, projectId.Id, clusterName).Execute()
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func validateCluster(cluster *admin.AdvancedClusterDescription, event MyEvent) bool {
	changed := false

	if event.ElectableSpecs != nil {
		for i := range *cluster.ReplicationSpecs {
			regionConfigs := *(*cluster.ReplicationSpecs)[i].RegionConfigs

			for j := range regionConfigs {
				if regionConfigs[j].ElectableSpecs != nil {
					if event.ElectableSpecs.InstanceSize != "" && event.ElectableSpecs.InstanceSize != regionConfigs[j].ElectableSpecs.GetInstanceSize() {
						log.Printf("Electable InstanceSize changed: %s -> %s", regionConfigs[j].ElectableSpecs.GetInstanceSize(), event.ElectableSpecs.InstanceSize)
						changed = true
					}
					if event.ElectableSpecs.DiskIOPS != 0 && event.ElectableSpecs.DiskIOPS != regionConfigs[j].ElectableSpecs.GetDiskIOPS() {
						log.Printf("Electable DiskIOPS changed: %d -> %d", regionConfigs[j].ElectableSpecs.GetDiskIOPS(), event.ElectableSpecs.DiskIOPS)
						changed = true
					}
					if event.ElectableSpecs.EBSVolumeType != "" && event.ElectableSpecs.EBSVolumeType != regionConfigs[j].ElectableSpecs.GetEbsVolumeType() {
						log.Printf("Electable EBSVolumeType changed: %s -> %s", regionConfigs[j].ElectableSpecs.GetEbsVolumeType(), event.ElectableSpecs.EBSVolumeType)
						changed = true
					}
				}
			}
		}
	}

	if event.ReadOnlySpecs != nil {
		for i := range *cluster.ReplicationSpecs {
			regionConfigs := *(*cluster.ReplicationSpecs)[i].RegionConfigs

			for j := range regionConfigs {
				if regionConfigs[j].ReadOnlySpecs != nil {
					if event.ReadOnlySpecs.InstanceSize != "" && event.ReadOnlySpecs.InstanceSize != regionConfigs[j].ReadOnlySpecs.GetInstanceSize() {
						log.Printf("ReadOnly InstanceSize changed: %s -> %s", regionConfigs[j].ReadOnlySpecs.GetInstanceSize(), event.ReadOnlySpecs.InstanceSize)
						changed = true
					}
					if event.ReadOnlySpecs.DiskIOPS != 0 && event.ReadOnlySpecs.DiskIOPS != regionConfigs[j].ReadOnlySpecs.GetDiskIOPS() {
						log.Printf("ReadOnly DiskIOPS changed: %d -> %d", regionConfigs[j].ReadOnlySpecs.GetDiskIOPS(), event.ReadOnlySpecs.DiskIOPS)
						changed = true
					}
					if event.ReadOnlySpecs.EBSVolumeType != "" && event.ReadOnlySpecs.EBSVolumeType != regionConfigs[j].ReadOnlySpecs.GetEbsVolumeType() {
						log.Printf("ReadOnly EBSVolumeType changed: %s -> %s", regionConfigs[j].ReadOnlySpecs.GetEbsVolumeType(), event.ReadOnlySpecs.EBSVolumeType)
						changed = true
					}
				}
			}
		}
	}

	if event.AnalyticsSpecs != nil {
		for i := range *cluster.ReplicationSpecs {
			regionConfigs := *(*cluster.ReplicationSpecs)[i].RegionConfigs

			for j := range regionConfigs {
				if regionConfigs[j].AnalyticsSpecs != nil {
					if event.AnalyticsSpecs.InstanceSize != "" && event.AnalyticsSpecs.InstanceSize != regionConfigs[j].AnalyticsSpecs.GetInstanceSize() {
						log.Printf("Analytics InstanceSize changed: %s -> %s", regionConfigs[j].AnalyticsSpecs.GetInstanceSize(), event.AnalyticsSpecs.InstanceSize)
						changed = true
					}
					if event.AnalyticsSpecs.DiskIOPS != 0 && event.AnalyticsSpecs.DiskIOPS != regionConfigs[j].AnalyticsSpecs.GetDiskIOPS() {
						log.Printf("Analytics DiskIOPS changed: %d -> %d", regionConfigs[j].AnalyticsSpecs.GetDiskIOPS(), event.AnalyticsSpecs.DiskIOPS)
						changed = true
					}
					if event.AnalyticsSpecs.EBSVolumeType != "" && event.AnalyticsSpecs.EBSVolumeType != regionConfigs[j].AnalyticsSpecs.GetEbsVolumeType() {
						log.Printf("Analytics EBSVolumeType changed: %s -> %s", regionConfigs[j].AnalyticsSpecs.GetEbsVolumeType(), event.AnalyticsSpecs.EBSVolumeType)
						changed = true
					}
				}
			}
		}
	}

	if !changed {
		return false
	}

	return true
}

func synchronizeInstanceSize(electable, readonly InstanceSizer) {
	if electable.GetInstanceSize() != "" {
		readonly.SetInstanceSize(electable.GetInstanceSize())
	}
	if electable.GetDiskIOPS() != 0 {
		readonly.SetDiskIOPS(electable.GetDiskIOPS())
	}
	if electable.GetEbsVolumeType() != "" {
		readonly.SetEbsVolumeType(electable.GetEbsVolumeType())
	}
}

func AutoScaling(ctx context.Context, project Project, event *MyEvent, sdk *admin.APIClient) error {

	cluster, err := GetCluster(ctx, project, event.Cluster, sdk)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	if !validateCluster(cluster, *event) {
		log.Println("No changes detected, exiting program.")
		os.Exit(0)
	}

	cluster.ConnectionStrings = nil
	for i := range *cluster.ReplicationSpecs {
		regionConfigs := *(*cluster.ReplicationSpecs)[i].RegionConfigs

		for j := range regionConfigs {
			if event.ElectableSpecs != nil {
				if event.ElectableSpecs.InstanceSize != "" {
					regionConfigs[j].ElectableSpecs.SetInstanceSize(event.ElectableSpecs.InstanceSize)
				}
				if event.ElectableSpecs.DiskIOPS != 0 {
					regionConfigs[j].ElectableSpecs.SetDiskIOPS(event.ElectableSpecs.DiskIOPS)
				}
				if event.ElectableSpecs.EBSVolumeType != "" {
					regionConfigs[j].ElectableSpecs.SetEbsVolumeType(event.ElectableSpecs.EBSVolumeType)
				}
				if event.ElectableSpecs.NodeCount != 0 {
					regionConfigs[j].ElectableSpecs.SetNodeCount(event.ElectableSpecs.NodeCount)
				}
				if event.ElectableSpecs.AutoScale.MaxInstanceSize != "" {
					regionConfigs[j].AutoScaling.Compute.SetMaxInstanceSize(event.ElectableSpecs.AutoScale.MaxInstanceSize)
				}
				if event.ElectableSpecs.AutoScale.MinInstanceSize != "" {
					regionConfigs[j].AutoScaling.Compute.SetMinInstanceSize(event.ElectableSpecs.AutoScale.MinInstanceSize)
				}
				synchronizeInstanceSize(event.ElectableSpecs, regionConfigs[j].ReadOnlySpecs)
				synchronizeInstanceSize(event.ElectableSpecs, regionConfigs[j].AnalyticsSpecs)
			} else {
				log.Println("event.ElectableSpecs is nil, skipping ElectableSpecs update")
			}

			if event.ReadOnlySpecs != nil {
				if event.ReadOnlySpecs.NodeCount != 0 {
					regionConfigs[j].ReadOnlySpecs.SetNodeCount(event.ReadOnlySpecs.NodeCount)
				}
			} else {
				log.Println("event.ReadOnlySpecs is nil, skipping ReadOnlySpecs update")
			}

			if event.AnalyticsSpecs != nil {
				if event.AnalyticsSpecs.InstanceSize != "" {
					regionConfigs[j].AnalyticsSpecs.SetInstanceSize(event.AnalyticsSpecs.InstanceSize)
				}
				if event.AnalyticsSpecs.NodeCount != 0 {
					regionConfigs[j].AnalyticsSpecs.SetNodeCount(event.AnalyticsSpecs.NodeCount)
				}
				if event.AnalyticsSpecs.AutoScale.MaxInstanceSize != "" {
					regionConfigs[j].AutoScaling.Compute.SetMaxInstanceSize(event.AnalyticsSpecs.AutoScale.MaxInstanceSize)
				}
				if event.AnalyticsSpecs.AutoScale.MinInstanceSize != "" {
					regionConfigs[j].AutoScaling.Compute.SetMinInstanceSize(event.AnalyticsSpecs.AutoScale.MinInstanceSize)
				}
			} else {
				log.Println("event.AnalyticsSpecs is nil, skipping AnalyticsSpecs update")
			}
		}

	}

	jsonBody, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling updated cluster: %v", err)
	}
	log.Println("Updated JSON Body:\n", string(jsonBody))

	_, response, err := sdk.ClustersApi.UpdateCluster(ctx, project.Id, cluster.GetName(), cluster).Execute()
	if err != nil {
		return fmt.Errorf("failed to update cluster: %w", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d, response: %s", response.StatusCode, response.Body)
	}

	return nil
}
