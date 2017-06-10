package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
)

func labelsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func dataSourceGoogleComputeSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSnapshotRead,

		Schema: map[string]*schema.Schema{
			//"filter": dataSourceFiltersSchema(),
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_disk_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// FIXME: sha256 ?
			"snapshot_encryption_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// FIXME: sha256 ?
			"source_disk_encryption_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_size_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"licenses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": labelsSchemaComputed(),
		},
	}
}

func dataSourceGoogleComputeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	labels := d.Get("labels").(map[string]interface{})
	log.Printf("[DEBUG] Labels %s", labels)


	if len(labels) > 0 {
		filter := ""
		log.Printf("[DEBUG] Labels length : %d", len(labels))
		for k, v := range labels {
			log.Printf("[DEBUG] Label key : '%s', value : '%s'", k, v)
			filter = fmt.Sprintf("%s (labels.%s eq %s)", filter, k, v)
		}
		log.Printf("[DEBUG] Labels filter : %s", filter)
	}

	snapshot, err := config.clientCompute.Snapshots.Get(
		project, d.Get("name").(string)).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore

			return fmt.Errorf("Snapshot Not Found : %s", d.Get("name"))
		}

		return fmt.Errorf("Error reading snapshot: %s", err)
	}
	d.Set("self_link", snapshot.SelfLink)
	d.Set("description", snapshot.Description)
	d.Set("snapshot_encryption_key", snapshot.SnapshotEncryptionKey)
	d.Set("source_disk_link", snapshot.SourceDisk)
	d.Set("source_disk_encryption_key", snapshot.SourceDiskEncryptionKey)
	d.Set("source_disk_id", snapshot.SourceDiskId)
	d.Set("status", snapshot.Status)
	d.Set("storage_size", snapshot.StorageBytes)
	d.Set("storage_size_status", snapshot.StorageBytesStatus)
	d.Set("disk_size", snapshot.DiskSizeGb)
	d.Set("labels", snapshot.Labels)

	d.SetId(snapshot.Name)
	return nil
}
