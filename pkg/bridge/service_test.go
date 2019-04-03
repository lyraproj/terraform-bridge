package bridge_test

import (
	"bytes"
	"testing"

	"github.com/lyraproj/pcore/types"
	"github.com/stretchr/testify/require"

	"github.com/terraform-providers/terraform-provider-github/github"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/terraform-bridge/pkg/bridge"
)

func TestCreateService_metadata(t *testing.T) {
	pcore.Do(func(c px.Context) {
		s := bridge.CreateService(c, github.Provider().(*schema.Provider), "TerraformGithub", "github",
			&terraform.ResourceConfig{
				Config: map[string]interface{}{},
			})
		_, md := s.Metadata(c)
		b := bytes.NewBufferString(``)
		mda := px.Wrap(c, md).(*types.Array)
		mda.ToString(b, px.Pretty, nil)
		mdr, err := types.Parse(b.String())
		require.Nil(t, err)
		mdb := types.ResolveDeferred(c, mdr, px.EmptyMap).(*types.Array)
		require.True(t, mda.Equals(mdb, nil))
	})
}

func TestCreateService_typeset(t *testing.T) {
	pcore.Do(func(c px.Context) {
		s := bridge.CreateService(c, github.Provider().(*schema.Provider), "TerraformGithub", "github",
			&terraform.ResourceConfig{
				Config: map[string]interface{}{},
			})
		ts, _ := s.Metadata(c)
		b := bytes.NewBufferString(``)
		ts.ToString(b, px.PrettyExpanded, nil)
		to := c.ParseType(b.String()).(px.ResolvableType).Resolve(c)
		require.True(t, ts.Equals(to, nil))
	})
}
