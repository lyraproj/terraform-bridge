package bridge

import (
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/lyraproj/pcore/px"
)

// GenericTfHandlerHandler ...
type GenericTFHandler struct {
	provider            *schema.Provider
	resourceType        px.Type
	externalIdAttribute string
	nativeType          string
}

func NewTFHandler(provider *schema.Provider, resourceType px.Type, externalIdAttributem, nativeType string) *GenericTFHandler {
	return &GenericTFHandler{provider, resourceType, strings.ToLower(externalIdAttributem), nativeType}
}

var once sync.Once
var Config *terraform.ResourceConfig

func configureProvider(p *schema.Provider) {
	once.Do(func() {
		if Config == nil {
			Config = &terraform.ResourceConfig{
				Config: map[string]interface{}{},
			}
		}
		err := p.Configure(Config)
		if err != nil {
			panic(err)
		}
	})
}

// Create ...
func (h *GenericTFHandler) Create(desired px.PuppetObject) (px.PuppetObject, string, error) {
	c := px.CurrentContext()
	log := hclog.Default()
	if log.IsInfo() {
		log.Info("Create "+h.nativeType, "desired", desired.String())
	}
	configureProvider(h.provider)
	rc := &terraform.ResourceConfig{
		Config: TerraformMarshal(c, desired),
	}
	id, err := Create(h.provider, h.nativeType, rc)
	if err != nil {
		return nil, "", err
	}
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
		log.Info("Update "+dt.Name(), "desired", desired.String())
	}
	configureProvider(h.provider)
	rc := &terraform.ResourceConfig{
		Config: TerraformMarshal(c, desired),
	}
	actual, err := Update(h.provider, h.nativeType, externalID, rc)
	if err != nil {
		return nil, err
	}
	x := TerraformUnMarshal(c, h.externalIdAttribute, externalID, actual, dt)
	if log.IsInfo() {
		log.Info("Update Actual State "+dt.Name(), "actual", x.String())
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
	configureProvider(h.provider)
	id, actual, err := Read(h.provider, h.nativeType, externalID)
	if err != nil {
		return nil, err
	}
	x := TerraformUnMarshal(c, h.externalIdAttribute, id, actual, h.resourceType.(px.ObjectType))
	if log.IsInfo() {
		log.Info("Read Actual State "+h.nativeType, "actual", x.String())
	}
	return x, nil
}

// Delete ...
func (h *GenericTFHandler) Delete(externalID string) error {
	log := hclog.Default()

	if log.IsInfo() {
		log.Info("Delete "+h.nativeType, "externalID", externalID)
	}
	configureProvider(h.provider)
	return Delete(h.provider, h.nativeType, externalID)
}
