package bridge

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/pcore/types"
)

// TerraformMarshal converts a PuppetObject into its Terraform representation
func TerraformMarshal(c px.Context, s px.PuppetObject) map[string]interface{} {
	t := s.PType().(px.ObjectType)

	attrs := t.AttributesInfo().Attributes()
	m := make(map[string]interface{}, len(attrs))

	for _, a := range attrs {
		tags := a.Tags(c)
		if tags != nil && tags.Tag(`lyra`) == ignoreTag {
			continue
		}

		pv := a.Get(s)
		if pv.Equals(px.Undef, nil) {
			continue
		}

		// Add the field value to the map
		marshaledValue := marshal(c, pv)
		if marshaledValue != nil {
			m[strings.ToLower(a.Name())] = marshaledValue
		}
	}
	return m
}

func marshal(c px.Context, v px.Value) interface{} {
	if v.Equals(px.Undef, nil) {
		return nil
	}
	switch v := v.(type) {
	case px.PuppetObject:
		return TerraformMarshal(c, v)
	case px.StringValue:
		return v.String()
	case px.Integer:
		return v.Int()
	case px.Boolean:
		return v.Bool()
	case px.Float:
		return v.Float()
	case *types.Timestamp:
		return v.Format(time.RFC3339)
	case *types.Regexp:
		return v.PatternString()
	case px.OrderedMap:
		nested := map[string]interface{}{}
		v.EachPair(func(k, v px.Value) { nested[k.String()] = marshal(c, v) })
		return nested
	case px.List:
		slice := make([]interface{}, v.Len())
		v.EachWithIndex(func(e px.Value, i int) { slice[i] = marshal(c, e) })
		return slice
	default:
		hclog.Default().Error(fmt.Sprintf("TerraformMarshal: Skipping unsupported pcore type: %s", v.PType()))
		return nil
	}
}

// TerraformUnMarshal converts a Terraform representation into a PuppetObject
func TerraformUnMarshal(c px.Context, extIdName, extId string, s map[string]interface{}, t px.ObjectType) px.PuppetObject {
	log := hclog.Default()
	if log.IsDebug() {
		log.Debug("TerraformUnMarshal", "type", t.String(), "state", px.Wrap(c, s).String())
	}
	attrs := t.AttributesInfo().Attributes()
	ie := make([]*types.HashEntry, 0, len(s))
	ie = append(ie, types.WrapHashEntry2(extIdName, types.WrapString(extId)))
	for _, a := range attrs {
		tfName := strings.ToLower(a.Name())
		if v, ok := s[tfName]; ok {
			ie = append(ie, types.WrapHashEntry2(a.Name(), types.CoerceTo(c, a.Label(), false, a.Type(), px.Wrap(c, v))))
		}
	}
	return px.New(c, t, types.WrapHash(ie)).(px.PuppetObject)
}
