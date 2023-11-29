package http

import (
	"errors"
	"fmt"

	"github.com/ChipArtem/k6/js/common"
	"github.com/dop251/goja"
)

//nolint:gochecknoglobals
var defaultExpectedStatuses = expectedStatuses{
	minmax: [][2]int{{200, 399}},
}

// expectedStatuses is specifically totally unexported so it can't be used for anything else but
// SetResponseCallback and nothing can be done from the js side to modify it or make an instance of
// it except using ExpectedStatuses
type expectedStatuses struct {
	minmax [][2]int
	exact  []int
}

func (e expectedStatuses) match(status int) bool {
	for _, v := range e.exact {
		if v == status {
			return true
		}
	}

	for _, v := range e.minmax {
		if v[0] <= status && status <= v[1] {
			return true
		}
	}
	return false
}

// expectedStatuses returns expectedStatuses object based on the provided arguments.
// The arguments must be either integers or object of `{min: <integer>, max: <integer>}`
// kind. The "integer"ness is checked by the Number.isInteger.
func (mi *ModuleInstance) expectedStatuses(args ...goja.Value) *expectedStatuses {
	rt := mi.vu.Runtime()

	if len(args) == 0 {
		common.Throw(rt, errors.New("no arguments"))
	}
	var result expectedStatuses

	jsIsInt, _ := goja.AssertFunction(rt.GlobalObject().Get("Number").ToObject(rt).Get("isInteger"))
	isInt := func(a goja.Value) bool {
		v, err := jsIsInt(goja.Undefined(), a)
		return err == nil && v.ToBoolean()
	}

	errMsg := "argument number %d to expectedStatuses was neither an integer nor an object like {min:100, max:329}"
	for i, arg := range args {
		o := arg.ToObject(rt)
		if o == nil {
			common.Throw(rt, fmt.Errorf(errMsg, i+1))
		}

		if isInt(arg) {
			result.exact = append(result.exact, int(o.ToInteger()))
		} else {
			min := o.Get("min")
			max := o.Get("max")
			if min == nil || max == nil {
				common.Throw(rt, fmt.Errorf(errMsg, i+1))
			}
			if !(isInt(min) && isInt(max)) {
				common.Throw(rt, fmt.Errorf("both min and max need to be integers for argument number %d", i+1))
			}

			result.minmax = append(result.minmax, [2]int{int(min.ToInteger()), int(max.ToInteger())})
		}
	}
	return &result
}

// SetResponseCallback sets the responseCallback to the value provided. Supported values are
// expectedStatuses object or a `null` which means that metrics shouldn't be tagged as failed and
// `http_req_failed` should not be emitted - the behaviour previous to this
func (c *Client) SetResponseCallback(val goja.Value) {
	if val != nil && !goja.IsNull(val) {
		// This is done this way as ExportTo exports functions to empty structs without an error
		if es, ok := val.Export().(*expectedStatuses); ok {
			c.responseCallback = es.match
		} else {
			common.Throw(
				c.moduleInstance.vu.Runtime(),
				fmt.Errorf("unsupported argument, expected http.expectedStatuses"),
			)
		}
	} else {
		c.responseCallback = nil
	}
}
