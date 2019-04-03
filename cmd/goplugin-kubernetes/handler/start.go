package handler

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/servicesdk/grpc"
	"github.com/lyraproj/terraform-bridge/pkg/bridge"
	"github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
)

const Namespace = `Kubernetes`

// Start this server running
func Start() {
	pcore.Do(func(c px.Context) {
		s := bridge.CreateService(c, kubernetes.Provider().(*schema.Provider), Namespace,
			&terraform.ResourceConfig{
				Config: map[string]interface{}{},
			})
		grpc.Serve(c, s)
	})
}
