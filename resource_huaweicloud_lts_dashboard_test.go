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

func getDadhboardResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, _ := httpclient_go.NewHttpClientGo(conf)
	c.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(conf, "lts", conf.Region, "v2/"+conf.HwClient.ProjectID+
		"/dashboards?id="+state.Primary.ID)
	response, err := c.Do()
	body, _ := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &entity.ReadDashBoardResp{}
	err = json.Unmarshal(body, rlt)
	if err != nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccessRule_basic(t *testing.T) {
	var instance entity.ReadDashBoardResp
	resourceName := "huaweicloud_lts_dashboard.dashboard_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDadhboardResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesDashBoard_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "log_group_id", "9ac33c09-7f00-4eed-b9a0-0ffaad7a64d1"),
					resource.TestCheckResourceAttr(resourceName, "log_group_name", "CTS"),
					resource.TestCheckResourceAttr(resourceName, "log_stream_id", "c3ab6968-a903-493d-a49a-5c45caaf32b4"),
					resource.TestCheckResourceAttr(resourceName, "log_stream_name", "test-znb"),
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

func tesDashBoard_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_dashboard" "dashboard_1" {
  log_group_id        = "9ac33c09-7f00-4eed-b9a0-0ffaad7a64d1"
  log_group_name      = "CTS"
  log_stream_id = "c3ab6968-a903-493d-a49a-5c45caaf32b4"
  log_stream_name = "test-znb"
  is_delete_charts      = "true"
  template_title   = ["cfw-log-analysis"]
  
}`)
}
