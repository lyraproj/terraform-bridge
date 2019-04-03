package bridge

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"

	"github.com/lyraproj/servicesdk/serviceapi"

	"github.com/lyraproj/servicesdk/service"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/pcore/types"
	"github.com/lyraproj/pcore/utils"

	// Ensure SDK is active
	_ "github.com/lyraproj/servicesdk/annotation"
)

// CreateService creates a serviceapi.Service from a Terraform schema.Provider
func CreateService(c px.Context, p *schema.Provider, ns string, config *terraform.ResourceConfig) serviceapi.Service {
	g := &generator{
		ctx:        c,
		provider:   &tfProvider{prv: p, cfg: config},
		generated:  make(map[string]bool, 53),
		skipPrefix: strings.ToLower(ns) + `_`,
		builder:    service.NewServiceBuilder(c, ns),
		namespace:  ns}

	g.builder.RegisterApiType(ns+`::GenericHandler`, NewTFHandler(nil, nil, ``, ``))
	b := bytes.NewBufferString(``)
	for k, v := range p.ResourcesMap {
		b.Reset()
		g.generateResourceType(k, v, b)
	}
	return g.builder.Server()
}

type generator struct {
	ctx        px.Context
	provider   *tfProvider
	generated  map[string]bool
	builder    *service.Builder
	namespace  string
	skipPrefix string
}

// isValidateFunc compare two SchemaValidateFunc pointers for equality
func isValidateFunc(actual, expected schema.SchemaValidateFunc) bool {
	return reflect.ValueOf(actual).Pointer() == reflect.ValueOf(expected).Pointer()
}

func (g *generator) writeRequiredType(s *schema.Schema, b *bytes.Buffer) {
	switch s.Type {
	case schema.TypeList, schema.TypeSet:
		// Terraform wraps nested resources in a TypeSet or TypeList because they don't have any
		// ValueType to represent the Object type. The intermediate array is elimitated here
		if r, ok := s.Elem.(*schema.Resource); ok && s.MaxItems == 1 {
			g.generateResourceType(``, r, b)
		} else {
			b.WriteString("Array[")
			g.writeElementType(s, b)
			minMax(s, b)
			b.WriteByte(']')
		}
	case schema.TypeMap:
		b.WriteString("Hash[String,")
		g.writeElementType(s, b)
		b.WriteByte(']')
	default:
		g.writePrimitive(s, s.Type, b)
	}
}

func (g *generator) writePrimitive(s *schema.Schema, t schema.ValueType, b *bytes.Buffer) {
	switch t {
	case schema.TypeBool:
		b.WriteString("Boolean")
	case schema.TypeInt:
		b.WriteString("Integer")
	case schema.TypeFloat:
		b.WriteString("Float")
	case schema.TypeString:
		if s != nil && s.ValidateFunc != nil {
			// Derive Timestamp and Regexp from ValidateFunc
			if isValidateFunc(s.ValidateFunc, validation.ValidateRFC3339TimeString) {
				b.WriteString("Timestamp")
				break
			} else if isValidateFunc(s.ValidateFunc, validation.ValidateRegexp) {
				b.WriteString("Regexp")
				break
			}
		}
		b.WriteString("String")
	default:
		panic(fmt.Errorf("not a primitive type: %s", t.String()))
	}
}

func minMax(s *schema.Schema, b *bytes.Buffer) {
	if s.MinItems != 0 {
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(s.MinItems))
		if s.MaxItems != 0 {
			b.WriteByte(',')
			b.WriteString(strconv.Itoa(s.MaxItems))
		}
	} else {
		if s.MaxItems != 0 {
			b.WriteString(`,0,`)
			b.WriteString(strconv.Itoa(s.MaxItems))
		}
	}
}

func (g *generator) writeElementType(s *schema.Schema, b *bytes.Buffer) {
	switch el := s.Elem.(type) {
	case nil:
		// Default to using a string
		b.WriteString("String")
	case *schema.Resource:
		g.generateResourceType(``, el, b)
	case *schema.Schema:
		g.writeType(el, b)
	case schema.ValueType:
		switch el {
		case schema.TypeList, schema.TypeSet:
			// No nested element type available. Default to string
			b.WriteString("Array[String")
			minMax(s, b)
			b.WriteByte(']')
		case schema.TypeMap:
			// No nested element type available. Default to string
			b.WriteString("Hash[String, String]")
		default:
			g.writePrimitive(nil, el, b)
		}
	default:
		panic(fmt.Sprintf("Unsupported type: %v", el))
	}
}

func isPcoreOptional(s *schema.Schema) bool {
	// An attribute with a specified Default doesn't have to be optional since its value will be
	// set at all times, either specifically or to the default (by pcore as well as by terraform).
	//
	// An attribute with a DefaultFunc must be considered optional since the function may return
	// different values at different times (environment settings etc. may affect what it returns).
	opt := s.Default == nil && (s.Computed || s.Optional || s.DefaultFunc != nil)

	// Terraform wraps nested resources in a TypeSet or TypeList because they don't have any
	// ValueType to represent the Object type. An s.MaxItems == 1 and s.MinItems == 0 for a
	// Resource therefore means that the entry is optional.
	if !opt && s.MaxItems == 1 && s.MinItems == 0 {
		_, opt = s.Elem.(*schema.Resource)
	}
	return opt
}

func (g *generator) writeType(s *schema.Schema, ab *bytes.Buffer) {
	opt := isPcoreOptional(s)
	if opt {
		ab.WriteString("Optional[")
	}
	g.writeRequiredType(s, ab)
	if opt {
		ab.WriteByte(']')
	}
}

func (g *generator) generateResourceType(nativeType string, r *schema.Resource, b *bytes.Buffer) {
	var provided []string
	var immutable []string

	// Determine field names and types
	b.WriteString(`Object{attributes=>{`)

	structType := ``
	extIdAttr := ``
	if nativeType != `` {
		lcStructType := g.skipPackage(nativeType)
		structType = strings.Title(lcStructType)

		// Attempt to use an external id that consists of the type name with an appended '_id'. If this
		// collides with an existing attribute, then use '_lyra_id' instead to ensure uniqueness.
		extIdSuffix := `_id`
		extIdAttr = lcStructType + extIdSuffix
		if _, exists := r.Schema[extIdAttr]; exists {
			extIdSuffix = `_lyra_id`
			extIdAttr = lcStructType + extIdSuffix
		}

		if _, ok := g.generated[nativeType]; ok {
			return
		}
		provided = append(provided, extIdAttr)
		utils.PuppetQuote(b, extIdAttr)
		b.WriteString(`=>Optional[String],`)
	}

	// Ensure names are sorted
	i := len(r.Schema)
	names := make([]string, i)
	for name := range r.Schema {
		i--
		names[i] = name
	}
	sort.Strings(names)

	first := true
	tb := bytes.NewBufferString(``)
	for _, name := range names {
		rs := r.Schema[name]
		if rs.Removed != `` || rs.Deprecated != `` {
			// TODO: Perhaps control inclusion of deprecated attributes using a flag
			continue
		}

		if first {
			first = false
		} else {
			b.WriteByte(',')
		}

		tb.Reset()
		g.writeType(rs, tb)
		ts := tb.String()

		utils.PuppetQuote(b, name)
		b.WriteString(`=>{type=>`)
		b.WriteString(ts)
		if rs.Default != nil {
			pt := g.ctx.ParseType(ts)
			b.WriteByte(',')
			g.writeDefaultValue(name+`.default`, rs.Default, pt, b)
		}
		b.WriteByte('}')
		if rs.ForceNew {
			immutable = append(immutable, name)
		}
		if rs.Computed {
			provided = append(provided, name)
		}
	}
	b.WriteByte('}') // End of attributes

	if len(immutable) > 0 || len(provided) > 0 {
		b.WriteString(`,annotations=>{Lyra::Resource=>{`)
		if len(provided) > 0 {
			writeQuotedArrayEntry(`providedAttributes`, provided, b)
		}
		if len(immutable) > 0 {
			if len(provided) > 0 {
				b.WriteByte(',')
			}
			writeQuotedArrayEntry(`immutableAttributes`, immutable, b)
		}
		b.WriteString(`}}`) // End of Annotation and annotations Hash
	}
	b.WriteByte('}') // End of Object type

	if nativeType != `` {
		rn := g.namespace + `::` + structType
		rt := px.NewNamedType(rn, b.String())
		g.builder.RegisterHandler(rn+`Handler`, NewTFHandler(g.provider, rt, extIdAttr, nativeType), rt)
	}
}

func writeQuotedArrayEntry(name string, values []string, b *bytes.Buffer) {
	b.WriteString(name)
	b.WriteString(`=>[`)
	for i, p := range values {
		if i > 0 {
			b.WriteByte(',')
		}
		utils.PuppetQuote(b, p)
	}
	b.WriteByte(']')
}

func (g *generator) writeDefaultValue(label string, dv interface{}, t px.Type, b *bytes.Buffer) {
	v := types.CoerceTo(g.ctx, label, t, px.Wrap(g.ctx, dv))
	b.WriteString(`value=>`)
	v.ToString(b, types.Program, nil)
}

func (g *generator) skipPackage(name string) string {
	if strings.HasPrefix(name, g.skipPrefix) {
		name = name[len(g.skipPrefix):]
	}
	return name
}
