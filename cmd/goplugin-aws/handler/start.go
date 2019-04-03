package handler

import (
	"sync"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/servicesdk/grpc"
	"github.com/lyraproj/terraform-bridge/pkg/bridge"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var once sync.Once

const Namespace = `Aws`

// Start this server running
func Start() {
	once.Do(func() {
		hclog.DefaultOptions = &hclog.LoggerOptions{
			Name:            Namespace,
			Level:           hclog.Debug,
			JSONFormat:      true,
			IncludeLocation: false,
		}
		issue.IncludeStacktrace(true)
	})

	pcore.Do(func(c px.Context) {
		s := bridge.CreateService(c, aws.Provider().(*schema.Provider), Namespace, &terraform.ResourceConfig{
			Config: map[string]interface{}{
				"region": "eu-west-1",
			},
		})
		grpc.Serve(c, s)
	})
}
