package bridge

import (
	"bytes"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/davecgh/go-spew/spew"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"github.com/stretchr/testify/require"
)

var original interface{}
var expected map[string]interface{}
var tfSchema map[string]*schema.Schema

type Kennel struct {
	Height *int
	Width  *int
	Depth  *int
}

type Dog struct {
	Size   int
	Colour string
	Home   Kennel
}

type Cat struct {
	Tail   int
	Colour *string
	Shape  *string
}

type Llama struct {
	Location *string
}

type Person struct {
	PersonID   string
	FirstName  string
	LastName   *string
	Cool       bool
	Cooler     bool
	Uncool     *bool
	Up         float64
	Down       *float64
	When       time.Time
	Match      *regexp.Regexp
	Age        int
	PetDog     Dog
	PetDoggy   *Dog
	NoDog      *Dog
	PetCats    []Cat
	PetMoggies *[]Cat
	PetFelines []*Cat
	PetKitties *[]*Cat
	PetLlama   *Llama
	Tags       map[string]string
	PTags      *map[string]string
	List       []string
	PList      *[]string
	PPList     *[]*string
	// ItemList  []Item
	// PItemList *[]Item
}

func init() {
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.SortKeys = true
	three := 3
	four := 4
	smith := "Smith"
	nothing := ""
	brownish := "brownish"
	brown := "brown"
	browny := "browny"
	tru := true
	fl := 5.678
	original = &Person{
		FirstName: "John",
		LastName:  &smith,
		Age:       23,
		Cool:      false,
		Cooler:    true,
		Uncool:    &tru,
		Up:        1.234,
		Down:      &fl,
		When:      time.Date(2019, 3, 22, 10, 53, 20, 0, time.UTC),
		Match:     regexp.MustCompile(`.*blue.*`),
		PetDog: Dog{
			12,
			"red",
			Kennel{
				Height: &three,
				Width:  &four,
			},
		},
		PetDoggy: &Dog{
			23,
			"yellow",
			Kennel{
				Height: &four,
				Width:  &three,
			},
		},
		PetCats: []Cat{
			{
				Tail:   15,
				Colour: &brownish,
			},
			{
				Tail:   16,
				Colour: &brown,
			},
		},
		PetMoggies: &[]Cat{
			{
				Tail:   17,
				Colour: &browny,
			},
			{
				Tail:   18,
				Colour: &brownish,
			},
		},
		PetFelines: []*Cat{
			{
				Tail:   27,
				Colour: &browny,
			},
			{
				Tail:   28,
				Colour: &brownish,
			},
		},
		PetKitties: &[]*Cat{
			{
				Tail:   19,
				Colour: &brown,
			},
			{
				Tail:   20,
				Colour: &browny,
			},
		},
		Tags:   map[string]string{"foo": "bar", "moo": "baa"},
		PTags:  &map[string]string{"foo2": "bar2", "moo2": "baa2"},
		List:   []string{"aa", "bb", "cc"},
		PList:  &[]string{"aa", "bb", "cc"},
		PPList: &[]*string{&smith, &nothing, &brownish},
		// ItemList:  []Item{Item{"aa"}, Item{"bb"}, Item{"cc"}},
		// PItemList: &[]Item{Item{"aa"}, Item{"bb"}, Item{"cc"}},
	}
	expected = map[string]interface{}{
		"firstName": "John",
		"lastName":  "Smith",
		"age":       23,
		"cool":      false,
		"cooler":    true,
		"uncool":    true,
		"up":        1.234,
		"down":      5.678,
		"when":      "2019-03-22T10:53:20Z",
		"match":     ".*blue.*",
		"petDog": []map[string]interface{}{{
			"colour": "red",
			"size":   12,
			"home": []map[string]interface{}{{
				"height": 3,
				"width":  4,
			}},
		}},
		"petDoggy": []map[string]interface{}{{
			"colour": "yellow",
			"size":   23,
			"home": []map[string]interface{}{{
				"height": 4,
				"width":  3,
			}},
		}},
		"petCats": []map[string]interface{}{
			{
				"tail":   15,
				"colour": "brownish",
			},
			{
				"tail":   16,
				"colour": "brown",
			},
		},
		"petMoggies": []map[string]interface{}{
			{
				"tail":   17,
				"colour": "browny",
			},
			{
				"tail":   18,
				"colour": "brownish",
			},
		},
		"petFelines": []map[string]interface{}{
			{
				"tail":   27,
				"colour": "browny",
			},
			{
				"tail":   28,
				"colour": "brownish",
			},
		},
		"petKitties": []map[string]interface{}{
			{
				"tail":   19,
				"colour": "brown",
			},
			{
				"tail":   20,
				"colour": "browny",
			},
		},
		"tags":   map[string]interface{}{"foo": "bar", "moo": "baa"},
		"pTags":  map[string]interface{}{"foo2": "bar2", "moo2": "baa2"},
		"list":   []interface{}{"aa", "bb", "cc"},
		"pList":  []interface{}{"aa", "bb", "cc"},
		"pPList": []interface{}{"Smith", "", "brownish"},
	}

	_string := &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	_string_opt := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	_int := &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	_int_opt := &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	_float := &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}
	_float_opt := &schema.Schema{
		Type:     schema.TypeFloat,
		Optional: true,
	}
	_bool := &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}

	_bool_opt := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true}

	_kennel := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"height": _int_opt,
			"width":  _int_opt,
			"depth":  _int_opt,
		}}

	_dog := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"colour": _string,
			"size":   _int,
			"home": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Elem:     _kennel,
			}}}

	_cat := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tail":   _int,
			"colour": _string,
			"shape":  _string,
		}}

	_llama := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location": _string,
		}}

	tfSchema = map[string]*schema.Schema{
		"firstName": _string,
		"lastName":  _string_opt,
		"cool":      _bool,
		"cooler":    _bool,
		"uncool":    _bool_opt,
		"up":        _float,
		"down":      _float_opt,
		"when":      _string,
		"match":     _string_opt,
		"age":       _int,
		"petDog": &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem:     _dog},
		"petDoggy": &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 0,
			MaxItems: 1,
			Elem:     _dog},
		"noDog": &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 0,
			MaxItems: 1,
			Elem:     _dog},
		"petCats": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     _cat},
		"petMoggies": &schema.Schema{
			Optional: true,
			Type:     schema.TypeList,
			Elem:     _cat},
		"petFelines": &schema.Schema{
			Required: true,
			MinItems: 0,
			Type:     schema.TypeList,
			Elem:     _cat},
		"petKitties": &schema.Schema{
			Optional: true,
			MinItems: 0,
			Type:     schema.TypeList,
			Elem:     _cat},
		"petLlama": &schema.Schema{
			MinItems: 0,
			MaxItems: 1,
			Type:     schema.TypeList,
			Elem:     _llama},
		"tags": &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem:     schema.TypeString},
		"pTags": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     schema.TypeString},
		"list": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     schema.TypeString},
		"pList": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     schema.TypeString},
		"pPList": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     schema.TypeString},
	}
}

func registerTypes(c px.Context) px.ObjectType {
	kt := c.Reflector().TypeFromReflect(`Kennel`, nil, reflect.TypeOf(&Kennel{}))
	dt := c.Reflector().TypeFromReflect(`Dog`, nil, reflect.TypeOf(&Dog{}))
	ct := c.Reflector().TypeFromReflect(`Cat`, nil, reflect.TypeOf(&Cat{}))
	lt := c.Reflector().TypeFromReflect(`Llama`, nil, reflect.TypeOf(&Llama{}))
	pt := c.Reflector().TypeFromReflect(`Person`, nil, reflect.TypeOf(&Person{}))
	px.AddTypes(c, kt, dt, ct, lt, pt)
	return pt
}

func TestTerraformMarshal(t *testing.T) {
	pcore.Do(func(c px.Context) {
		registerTypes(c)
		actual := TerraformMarshal(c, px.Wrap(c, original).(px.PuppetObject), tfSchema)

		s1 := bytes.NewBufferString(``)
		px.Wrap(c, expected).ToString(s1, px.Pretty, nil)
		s2 := bytes.NewBufferString(``)
		px.Wrap(c, actual).ToString(s2, px.Pretty, nil)

		require.Equal(t, s1.String(), s2.String())
	})
}

func TestTerraformUnmarshal(t *testing.T) {
	pcore.Do(func(c px.Context) {
		pt := registerTypes(c)
		org := px.Wrap(c, original)
		actual := TerraformUnMarshal(c, `personID`, ``, expected, pt)
		s1 := bytes.NewBufferString(``)
		org.ToString(s1, px.Pretty, nil)
		s2 := bytes.NewBufferString(``)
		actual.ToString(s2, px.Pretty, nil)
		require.Equal(t, s1.String(), s2.String())
		require.True(t, org.Equals(actual, nil))
	})
}
