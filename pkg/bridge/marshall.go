package bridge

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/pcore/types"
)

// TerraformMarshal converts a PuppetObject into its Terraform representation
func TerraformMarshal(c px.Context, s px.PuppetObject, ts map[string]*schema.Schema) map[string]interface{} {
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
		dn := strings.ToLower(a.Name())
		marshaledValue := marshal(c, pv, ts[dn])
		if marshaledValue != nil {
			m[dn] = marshaledValue
		}
	}
	return m
}

var noSchemas = map[string]*schema.Schema{}

func marshal(c px.Context, v px.Value, ts *schema.Schema) interface{} {
	if v.Equals(px.Undef, nil) {
		return nil
	}
	switch v := v.(type) {
	case px.PuppetObject:
		if ts != nil {
			switch ts.Type {
			case schema.TypeList, schema.TypeSet:
				schemas := noSchemas
				if rs, ok := ts.Elem.(*schema.Resource); ok {
					schemas = rs.Schema
				}
				return []map[string]interface{}{TerraformMarshal(c, v, schemas)}
			}
		}
		return TerraformMarshal(c, v, noSchemas)
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
		var es *schema.Schema
		nested := map[string]interface{}{}
		if ts != nil && ts.Type == schema.TypeMap {
			if s, ok := ts.Elem.(*schema.Schema); ok {
				es = s
			}
		}
		v.EachPair(func(k, v px.Value) {
			nested[k.String()] = marshal(c, v, es)
		})
		return nested
	case px.List:
		var es *schema.Schema
		slice := make([]interface{}, v.Len())
		if ts != nil && (ts.Type == schema.TypeList || ts.Type == schema.TypeSet) {
			if s, ok := ts.Elem.(*schema.Schema); ok {
				es = s
			}
		}
		v.EachWithIndex(func(e px.Value, i int) { slice[i] = marshal(c, e, es) })
		return slice
	default:
		hclog.Default().Error(fmt.Sprintf("TerraformMarshal: Skipping unsupported pcore type: %s", v.PType()))
		return nil
	}
}

// TerraformUnMarshal converts a Terraform representation into a PuppetObject
func TerraformUnMarshal(c px.Context, extIdName, extId string, s map[string]interface{}, t px.ObjectType) px.PuppetObject {
	h := px.Wrap(c, s).(px.OrderedMap)

	log := hclog.Default()
	if log.IsDebug() {
		b := bytes.NewBufferString(``)
		h.ToString(b, px.Pretty, nil)
		log.Debug("TerraformUnMarshal before conversion", "type", t.String(), "state", b.String())
	}

	attrs := t.AttributesInfo().Attributes()
	h = optObjToOneElementArray(h, attrs)

	if log.IsDebug() {
		b := bytes.NewBufferString(``)
		h.ToString(b, px.Pretty, nil)
		log.Debug("TerraformUnMarshal after conversion", "type", t.String(), "state", b.String())
	}

	ie := make([]*types.HashEntry, 0, len(s))
	ie = append(ie, types.WrapHashEntry2(extIdName, types.WrapString(extId)))
	for _, a := range attrs {
		if v, ok := h.Get4(a.Name()); ok {
			ie = append(ie, types.WrapHashEntry2(a.Name(), types.CoerceTo(c, a.Label(), false, a.Type(), v)))
		}
	}
	return px.New(c, t, types.WrapHash(ie)).(px.PuppetObject)
}

// Terraform will store an Object as one-element array to be able to declare the schema in the Elem
// of that array.
func optObjToOneElementArray(rMap px.OrderedMap, attrs []px.Attribute) px.OrderedMap {
	es := make([]*types.HashEntry, 0, rMap.Len())
	for _, a := range attrs {
		n := types.WrapString(a.Name())
		if v, ok := rMap.Get(n); ok {
			es = append(es, types.WrapHashEntry(n, convertValue(v, a.Type())))
		}
	}
	return types.WrapHash(es)
}

func convertValue(v px.Value, t px.Type) px.Value {
	switch at := t.(type) {
	case *types.OptionalType:
		v = convertValue(v, at.ContainedType())
	case px.ObjectType:
		if a, ok := v.(*types.Array); ok && a.Len() == 1 {
			v = a.At(0)
		}
		if h, ok := v.(*types.Hash); ok {
			v = optObjToOneElementArray(h, at.AttributesInfo().Attributes())
		}
	case *types.StructType:
		if a, ok := v.(*types.Array); ok && a.Len() == 1 {
			v = a.At(0)
		}
		ts := at.Elements()
		if h, ok := v.(*types.Hash); ok {
			es := make([]*types.HashEntry, 0, h.Len())
			for _, se := range ts {
				n := types.WrapString(se.Name())
				if v, ok := h.Get(n); ok {
					es = append(es, types.WrapHashEntry(n, convertValue(v, se.Value())))
				}
			}
			v = types.WrapHash(es)
		}
	case *types.HashType:
		if h, ok := v.(*types.Hash); ok {
			vt := at.ValueType()
			v = h.MapValues(func(v px.Value) px.Value { return convertValue(v, vt) })
		}
	case *types.ArrayType:
		if h, ok := v.(*types.Array); ok {
			et := at.ElementType()
			v = h.Map(func(v px.Value) px.Value { return convertValue(v, et) })
		}
	}
	return v
}
