package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ppreeper/odoojrpc"
)

var mulPatternTests = []struct {
	filter    string
	rexpected odoojrpc.FilterArg
}{
	{"[]", nil},
	{"[[]]", nil},
	{"[)", nil},
	{"[)]", nil},
	{"[()]", nil},
	{"()", nil},
	{"[('s','=','A')]",
		odoojrpc.FilterArg{odoojrpc.FilterArg{"s", "=", "A"}}},
	{"('c','=',3)", nil},
	{"[('c','=',3]", nil},
	{"[('c','=',3)]",
		odoojrpc.FilterArg{odoojrpc.FilterArg{"c", "=", 3}}},
	{"[('c','=',3.3)]",
		odoojrpc.FilterArg{odoojrpc.FilterArg{"c", "=", 3.3}}},
	{"['!',('purchase_ok','=',True)]",
		odoojrpc.FilterArg{"!", odoojrpc.FilterArg{"purchase_ok", "=", true}}},
	{"[('c','=',3),('s','=','A')]",
		odoojrpc.FilterArg{odoojrpc.FilterArg{"c", "=", 3}, odoojrpc.FilterArg{"s", "=", "A"}}},
	{"[('purchase_ok','=',True),('sale_ok','=',True)]",
		odoojrpc.FilterArg{odoojrpc.FilterArg{"purchase_ok", "=", true}, odoojrpc.FilterArg{"sale_ok", "=", true}}},
	{"['|',('purchase_ok','=',True),('sale_ok','=',True)]",
		odoojrpc.FilterArg{"|", odoojrpc.FilterArg{"purchase_ok", "=", true}, odoojrpc.FilterArg{"sale_ok", "=", true}}},
	{"['|','&',('purchase_ok','=',True),('sale_ok','=',True),'&',('landed_cost_ok','=',True),('type','=','service')]",
		odoojrpc.FilterArg{
			"|",
			"&",
			odoojrpc.FilterArg{"purchase_ok", "=", true},
			odoojrpc.FilterArg{"sale_ok", "=", true},
			"&",
			odoojrpc.FilterArg{"landed_cost_ok", "=", true},
			odoojrpc.FilterArg{"type", "=", "service"},
		}},
	{"['&','&',('state','in',('sale','done')),('is_service','=',True),'|',('project_id','!=',False),('task_id','!=',False)]",
		odoojrpc.FilterArg{
			"&",
			"&",
			odoojrpc.FilterArg{"state", "in", odoojrpc.FilterArg{"sale", "done"}},
			odoojrpc.FilterArg{"is_service", "=", true},
			"|",
			odoojrpc.FilterArg{"project_id", "!=", false},
			odoojrpc.FilterArg{"task_id", "!=", false},
		}},
}

func TestParseFilter(t *testing.T) {
	for _, mt := range mulPatternTests {
		rfilter, rerror := parseFilter(mt.filter)
		fmt.Printf("\n----test----\nfilter:\t%v\nexpect:\t%v\ngot:\t%v\nerror:\t%v\nequal:\t%v\n", mt.filter, mt.rexpected, rfilter, rerror, reflect.DeepEqual(mt.rexpected, rfilter))
		if !reflect.DeepEqual(mt.rexpected, rfilter) {
			t.Errorf("no match: expected: %v T:%v got: %v T:%v", mt.rexpected, reflect.TypeOf(mt.rexpected), rfilter, reflect.TypeOf(rfilter))
		}
	}
}
