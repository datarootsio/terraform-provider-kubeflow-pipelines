package kubeflowpipelines

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceKubeflowPipelinesPipeline_basic(t *testing.T) {
	resourceName := "data.kubeflowpipelines_pipeline.test"
	pipelineName := acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubeflowPipelinesPipelineBasic(pipelineName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", pipelineName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("Description %s", pipelineName)),
				),
			},
		},
	})
}

func testAccDataSourceKubeflowPipelinesPipelineBasic(pipelineName string) string {
	f, _ := os.Open("./tests/pipeline.yaml")
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return fmt.Sprintf(`
resource "kubeflowpipelines_pipeline" "test" {
  name          = "%s"
  description   = "Description %s"
  file_base64   = "%s"
}

data "kubeflowpipelines_pipeline" "test" {
  id = kubeflowpipelines_pipeline.test.id
}
`, pipelineName, pipelineName, encoded)
}