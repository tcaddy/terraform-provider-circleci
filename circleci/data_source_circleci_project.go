package circleci

import (
	"fmt"

	client "github.com/SectorLabs/terraform-provider-circleci/circleci/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCircleCIProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCircleCIProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the project",
			},
		},
	}
}

func dataSourceCircleCIProjectRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	name := d.Get("name").(string)

	project, err := c.GetProject(name)
	if err != nil {
		return fmt.Errorf("failed call to GetProject: %w", err)
	}

	d.SetId(project.ID)
	return nil
}
