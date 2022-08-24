package cmdb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

func getStructTemplateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "lts", conf.Region, "v2/"+conf.HwClient.ProjectID+
		"/lts/struct/template?logGroupId="+state.Primary.ID)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.CreateLogtankResponse{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestLtsStructTemplate_basic(t *testing.T) {
	var instance entity.CreateLogtankResponse
	resourceName := "huaweicloud_lts_struct_template.struct_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getStructTemplateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesLtsStructTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "log_group_id", "e5e5a6de-354d-45af-a9da-0fb91e9a3796"),
					resource.TestCheckResourceAttr(resourceName, "log_topic_id", "45bbeee7-2144-4d40-9c80-ba452298b6b8"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func tesLtsStructTemplate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_elb_log" "elb_1" {
  log_group_id        = "e5e5a6de-354d-45af-a9da-0fb91e9a3796"
  log_topic_id      = "45bbeee7-2144-4d40-9c80-ba452298b6b8"
  template_type    = "custom"
}`)
}
