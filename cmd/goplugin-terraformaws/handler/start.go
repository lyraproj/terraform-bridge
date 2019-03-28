package handler

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/servicesdk/grpc"
	"github.com/lyraproj/servicesdk/service"
	gp "github.com/lyraproj/terraform-bridge/cmd/goplugin-terraformaws/generated"
	"github.com/lyraproj/terraform-bridge/pkg/bridge"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var once sync.Once

// Server configures the Terraform provider and creates an instance of the server
func Server(c px.Context) *service.Server {
	sb := service.NewServiceBuilder(c, "TerraformAws")
	gp.Initialize(sb, aws.Provider().(*schema.Provider))
	return sb.Server()
}

// Start this server running
func Start() {
	once.Do(func() {
		hclog.DefaultOptions = &hclog.LoggerOptions{
			Name:            "TerraformAws",
			Level:           hclog.Debug,
			JSONFormat:      true,
			IncludeLocation: false,
		}
		issue.IncludeStacktrace(true)
	})

	bridge.Config = &terraform.ResourceConfig{
		Config: map[string]interface{}{
			"region": "eu-west-1",
		},
	}
	pcore.Do(func(c px.Context) {
		grpc.Serve(c, Server(c))
	})
}
