package kubeflowpipelines

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubeflow/pipelines/backend/api/go_http_client/experiment_client/experiment_service"
)

func dataSourceKubeflowPipelinesExperiment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubeflowPipelinesExperimentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubeflowPipelinesExperimentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Meta).Experiment
	context := meta.(*Meta).Context

	name := d.Get("name").(string)
	id := d.Get("id").(string)

	if id == "" {
		resp, err := client.ExperimentService.ListExperiment(nil, nil)
		if err != nil {
			return fmt.Errorf("unable to get list of experiments: %s", name)
		}

		experimentFound := false

		for _, item := range resp.Payload.Experiments {
			if item.Name == name {
				d.SetId(item.ID)
				d.Set("name", item.Name)
				d.Set("description", item.Description)
				experimentFound = true
				break
			}
		}

		if !experimentFound {
			return fmt.Errorf("unable to get experiment: %s", name)
		}
	} else {
		experimentParams := experiment_service.GetExperimentParams{
			ID:      id,
			Context: context,
		}

		resp, err := client.ExperimentService.GetExperiment(&experimentParams, nil)
		if err != nil {
			return fmt.Errorf("unable to get experiment: %s", id)
		}
		d.SetId(resp.Payload.ID)
		d.Set("name", resp.Payload.Name)
		d.Set("description", resp.Payload.Description)
	}
	return nil
}
