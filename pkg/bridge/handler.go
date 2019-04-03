package bridge

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/flatmap"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/pcore/px"
)

// GenericTfHandlerHandler ...
type GenericTFHandler struct {
	tfProvider          *tfProvider
	resourceType        px.Type
	externalIdAttribute string
	nativeType          string
}

type tfProvider struct {
	prv  *schema.Provider
	cfg  *terraform.ResourceConfig
	once sync.Once
}

func (p *tfProvider) provider() *schema.Provider {
	p.once.Do(func() {
		err := p.prv.Configure(p.cfg)
		if err != nil {
			panic(err)
		}
	})
	return p.prv
}

func NewTFHandler(provider *tfProvider, resourceType px.Type, externalIdAttributem, nativeType string) *GenericTFHandler {
	return &GenericTFHandler{provider, resourceType, strings.ToLower(externalIdAttributem), nativeType}
}

func (h *GenericTFHandler) provider() *schema.Provider {
	return h.tfProvider.provider()
}

func (h *GenericTFHandler) resource() *schema.Resource {
	if r, ok := h.provider().ResourcesMap[h.nativeType]; ok {
		return r
	}
	panic(fmt.Errorf("unknown resource type: %s", h.nativeType))
}

// Create ...
func (h *GenericTFHandler) Create(desired px.PuppetObject) (px.PuppetObject, string, error) {
	c := px.CurrentContext()
	log := hclog.Default()
	if log.IsInfo() {
		b := bytes.NewBufferString(``)
		desired.ToString(b, px.Pretty, nil)
		log.Info("Create "+h.nativeType, "desired", b.String())
	}
	rc := &terraform.ResourceConfig{
		Config: TerraformMarshal(c, desired, h.resource().Schema),
	}

	// To get Terraform to create a new resource, the ID must be blank and existing state must be empty (since the
	// resource does not exist yet), and the diff object should have no old state and all of the new state.
	info := &terraform.InstanceInfo{Type: h.nativeType}
	state := &terraform.InstanceState{
		Attributes: map[string]string{},
		Meta:       map[string]interface{}{},
	}
	prv := h.provider()
	diff, err := prv.Diff(info, state, rc)
	if err != nil {
		return nil, "", err
	}
	state, err = prv.Apply(info, state, diff)
	if err != nil {
		return nil, "", err
	}
	id := state.ID
	actual, err := h.Read(id)
	if err != nil {
		return nil, "", err
	}
	return actual, id, nil
}

// Update ...
func (h *GenericTFHandler) Update(externalID string, desired px.PuppetObject) (px.PuppetObject, error) {
	c := px.CurrentContext()
	log := hclog.Default()
	dt := desired.PType().(px.ObjectType)
	if log.IsInfo() {
		b := bytes.NewBufferString(``)
		desired.ToString(b, px.Pretty, nil)
		log.Info("Update "+dt.Name(), "desired", b.String())
	}

	rc := &terraform.ResourceConfig{
		Config: TerraformMarshal(c, desired, h.resource().Schema),
	}

	info := &terraform.InstanceInfo{Type: h.nativeType}
	state := &terraform.InstanceState{
		ID:         externalID,
		Attributes: map[string]string{},
		Meta:       map[string]interface{}{},
	}
	prv := h.provider()
	state, err := prv.Refresh(info, state)
	if err != nil {
		return nil, err
	}
	diff, err := prv.Diff(info, state, rc)
	if err != nil {
		return nil, err
	}
	if diff == nil {
		hclog.Default().Debug("Update", "type", h.nativeType, "msg", "diff is zero")
	} else {
		state, err = prv.Apply(info, state, diff)
		if err != nil {
			return nil, err
		}
	}
	actual := expand(state)
	x := TerraformUnMarshal(c, h.externalIdAttribute, externalID, actual, dt)
	if log.IsDebug() {
		b := bytes.NewBufferString(``)
		x.ToString(b, px.Pretty, nil)
		log.Debug("Update Actual State "+dt.Name(), "actual", x.String())
	}
	return x, nil
}

// Read ...
func (h *GenericTFHandler) Read(externalID string) (px.PuppetObject, error) {
	c := px.CurrentContext()
	log := hclog.Default()

	if log.IsInfo() {
		log.Info("Read "+h.nativeType, "externalID", externalID)
	}

	info := &terraform.InstanceInfo{Type: h.nativeType}
	state := &terraform.InstanceState{
		ID:         externalID,
		Attributes: map[string]string{},
		Meta:       map[string]interface{}{},
	}
	state, err := h.provider().Refresh(info, state)
	if err != nil {
		return nil, err
	}
	actual := expand(state)
	x := TerraformUnMarshal(c, h.externalIdAttribute, externalID, actual, h.resourceType.(px.ObjectType))
	if log.IsInfo() {
		b := bytes.NewBufferString(``)
		x.ToString(b, px.Pretty, nil)
		log.Info("Read Actual State "+h.nativeType, "actual", b.String())
	}
	return x, nil
}

// Delete ...
func (h *GenericTFHandler) Delete(externalID string) error {
	log := hclog.Default()
	if log.IsInfo() {
		log.Info("Delete "+h.nativeType, "externalID", externalID)
	}

	info := &terraform.InstanceInfo{Type: h.nativeType}
	state := &terraform.InstanceState{ID: externalID}
	diff := &terraform.InstanceDiff{Destroy: true}
	_, err := h.provider().Apply(info, state, diff)
	return err
}

func expand(state *terraform.InstanceState) map[string]interface{} {
	var outs map[string]interface{}
	if state != nil {
		outs = make(map[string]interface{})
		attrs := state.Attributes
		for _, key := range flatmap.Map(attrs).Keys() {
			outs[key] = flatmap.Expand(attrs, key)
		}
	}
	return outs
}
