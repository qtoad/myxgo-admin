package types

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/qtoad/mygo-admin/modules/config"
	"github.com/qtoad/mygo-admin/template/types/form"
)

type DisplayFnGenerator interface {
	Get(args ...interface{}) FieldFilterFunc
	JS() template.HTML
	HTML() template.HTML
}

type BaseDisplayFnGenerator struct{}

func (base *BaseDisplayFnGenerator) JS() template.HTML   { return "" }
func (base *BaseDisplayFnGenerator) HTML() template.HTML { return "" }

var displayFnGens = make(map[string]DisplayFnGenerator)

func RegisterDisplayFnGenerator(key string, gen DisplayFnGenerator) {
	if _, ok := displayFnGens[key]; ok {
		panic("display function generator has been registered")
	}
	displayFnGens[key] = gen
}

type FieldDisplay struct {
	Display              FieldFilterFunc
	DisplayProcessChains DisplayProcessFuncChains
}

func (f FieldDisplay) ToDisplay(value FieldModel) interface{} {
	val := f.Display(value)

	if len(f.DisplayProcessChains) > 0 && f.IsNotSelectRes(val) {
		valStr := fmt.Sprintf("%v", val)
		for _, process := range f.DisplayProcessChains {
			valStr = fmt.Sprintf("%v", process(FieldModel{
				Row:   value.Row,
				Value: valStr,
				ID:    value.ID,
			}))
		}
		return valStr
	}

	return val
}

func (f FieldDisplay) IsNotSelectRes(v interface{}) bool {
	switch v.(type) {
	case template.HTML:
		return false
	case []string:
		return false
	case [][]string:
		return false
	default:
		return true
	}
}

func (f FieldDisplay) ToDisplayHTML(value FieldModel) template.HTML {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return h
	} else if s, ok := v.(string); ok {
		return template.HTML(s)
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return template.HTML(arr[0])
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		return arr[0]
	} else if v != nil {
		return template.HTML(fmt.Sprintf("%v", v))
	} else {
		return ""
	}
}

func (f FieldDisplay) ToDisplayString(value FieldModel) string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return string(h)
	} else if s, ok := v.(string); ok {
		return s
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return arr[0]
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		return string(arr[0])
	} else if v != nil {
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

func (f FieldDisplay) ToDisplayStringArray(value FieldModel) []string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return []string{string(h)}
	} else if s, ok := v.(string); ok {
		return []string{s}
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return arr
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		ss := make([]string, len(arr))
		for k, a := range arr {
			ss[k] = string(a)
		}
		return ss
	} else if v != nil {
		return []string{fmt.Sprintf("%v", v)}
	} else {
		return []string{}
	}
}

func (f FieldDisplay) ToDisplayStringArrayArray(value FieldModel) [][]string {
	v := f.ToDisplay(value)
	if h, ok := v.(template.HTML); ok {
		return [][]string{{string(h)}}
	} else if s, ok := v.(string); ok {
		return [][]string{{s}}
	} else if arr, ok := v.([]string); ok && len(arr) > 0 {
		return [][]string{arr}
	} else if arr, ok := v.([][]string); ok && len(arr) > 0 {
		return arr
	} else if arr, ok := v.([]template.HTML); ok && len(arr) > 0 {
		ss := make([]string, len(arr))
		for k, a := range arr {
			ss[k] = string(a)
		}
		return [][]string{ss}
	} else if v != nil {
		return [][]string{{fmt.Sprintf("%v", v)}}
	} else {
		return [][]string{}
	}
}

func (f FieldDisplay) AddLimit(limit int) DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		if limit > len(value.Value) {
			return value
		} else if limit < 0 {
			return ""
		} else {
			return value.Value[:limit]
		}
	})
}

func (f FieldDisplay) AddTrimSpace() DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.TrimSpace(value.Value)
	})
}

func (f FieldDisplay) AddSubstr(start int, end int) DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		if start > end || start > len(value.Value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value.Value) {
			end = len(value.Value)
		}
		return value.Value[start:end]
	})
}

func (f FieldDisplay) AddToTitle() DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.Title(value.Value)
	})
}

func (f FieldDisplay) AddToUpper() DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.ToUpper(value.Value)
	})
}

func (f FieldDisplay) AddToLower() DisplayProcessFuncChains {
	return f.DisplayProcessChains.Add(func(value FieldModel) interface{} {
		return strings.ToLower(value.Value)
	})
}

type DisplayProcessFuncChains []FieldFilterFunc

func (d DisplayProcessFuncChains) Valid() bool {
	return len(d) > 0
}

func (d DisplayProcessFuncChains) Add(f FieldFilterFunc) DisplayProcessFuncChains {
	return append(d, f)
}

func (d DisplayProcessFuncChains) Append(f DisplayProcessFuncChains) DisplayProcessFuncChains {
	return append(d, f...)
}

func (d DisplayProcessFuncChains) Copy() DisplayProcessFuncChains {
	if len(d) == 0 {
		return make(DisplayProcessFuncChains, 0)
	} else {
		var newDisplayProcessFnChains = make(DisplayProcessFuncChains, len(d))
		copy(newDisplayProcessFnChains, d)
		return newDisplayProcessFnChains
	}
}

func chooseDisplayProcessChains(internal DisplayProcessFuncChains) DisplayProcessFuncChains {
	if len(internal) > 0 {
		return internal
	}
	return globalDisplayProcessChains.Copy()
}

var globalDisplayProcessChains = make(DisplayProcessFuncChains, 0)

func AddGlobalDisplayProcessFunc(filterFunc FieldFilterFunc) {
	globalDisplayProcessChains = globalDisplayProcessChains.Add(filterFunc)
}

func AddLimit(limit int) DisplayProcessFuncChains {
	return addLimit(limit, globalDisplayProcessChains)
}

func AddTrimSpace() DisplayProcessFuncChains {
	return addTrimSpace(globalDisplayProcessChains)
}

func AddSubstr(start int, end int) DisplayProcessFuncChains {
	return addSubstr(start, end, globalDisplayProcessChains)
}

func AddToTitle() DisplayProcessFuncChains {
	return addToTitle(globalDisplayProcessChains)
}

func AddToUpper() DisplayProcessFuncChains {
	return addToUpper(globalDisplayProcessChains)
}

func AddToLower() DisplayProcessFuncChains {
	return addToLower(globalDisplayProcessChains)
}

func AddXssFilter() DisplayProcessFuncChains {
	return addXssFilter(globalDisplayProcessChains)
}

func AddXssJsFilter() DisplayProcessFuncChains {
	return addXssJsFilter(globalDisplayProcessChains)
}

func addLimit(limit int, chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		if limit > len(value.Value) {
			return value
		} else if limit < 0 {
			return ""
		} else {
			return value.Value[:limit]
		}
	})
	return chains
}

func addTrimSpace(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.TrimSpace(value.Value)
	})
	return chains
}

func addSubstr(start int, end int, chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		if start > end || start > len(value.Value) || end < 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(value.Value) {
			end = len(value.Value)
		}
		return value.Value[start:end]
	})
	return chains
}

func addToTitle(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.Title(value.Value)
	})
	return chains
}

func addToUpper(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.ToUpper(value.Value)
	})
	return chains
}

func addToLower(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return strings.ToLower(value.Value)
	})
	return chains
}

func addXssFilter(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		return html.EscapeString(value.Value)
	})
	return chains
}

func addXssJsFilter(chains DisplayProcessFuncChains) DisplayProcessFuncChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		replacer := strings.NewReplacer("<script>", "&lt;script&gt;", "</script>", "&lt;/script&gt;")
		return replacer.Replace(value.Value)
	})
	return chains
}

func setDefaultDisplayFuncOfFormType(f *FormPanel, typ form.Type) {
	if typ.IsMultiFile() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			if value.Value == "" {
				return ""
			}
			arr := strings.Split(value.Value, ",")
			res := "["
			for i, item := range arr {
				if i == len(arr)-1 {
					res += "'" + config.GetStore().URL(item) + "']"
				} else {
					res += "'" + config.GetStore().URL(item) + "',"
				}
			}
			return res
		}
	}
	if typ.IsSelect() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			return strings.Split(value.Value, ",")
		}
	}
}
