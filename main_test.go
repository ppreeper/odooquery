package main

import (
	"fmt"
	"reflect"
	"testing"
)

var mulPatternTests = []struct {
	filter    string
	rexpected []any
}{
	{"[]", nil},
	{"[[]]", nil},
	{"[)", nil},
	{"[)]", nil},
	{"[()]", nil},
	{"()", nil},
	{
		"[('s','=','A')]",
		[]any{[]any{"s", "=", "A"}},
	},
	{"('c','=',3)", nil},
	{"[('c','=',3]", nil},
	{
		"[('c','=',3)]",
		[]any{[]any{"c", "=", 3}},
	},
	{
		"[('c','=',3.3)]",
		[]any{[]any{"c", "=", 3.3}},
	},
	{
		"['!',('purchase_ok','=',True)]",
		[]any{"!", []any{"purchase_ok", "=", true}},
	},
	{
		"[('c','=',3),('s','=','A')]",
		[]any{[]any{"c", "=", 3}, []any{"s", "=", "A"}},
	},
	{
		"[('purchase_ok','=',True),('sale_ok','=',True)]",
		[]any{[]any{"purchase_ok", "=", true}, []any{"sale_ok", "=", true}},
	},
	{
		"['|',('purchase_ok','=',True),('sale_ok','=',True)]",
		[]any{"|", []any{"purchase_ok", "=", true}, []any{"sale_ok", "=", true}},
	},
	{
		"['|','&',('purchase_ok','=',True),('sale_ok','=',True),'&',('landed_cost_ok','=',True),('type','=','service')]",
		[]any{
			"|",
			"&",
			[]any{"purchase_ok", "=", true},
			[]any{"sale_ok", "=", true},
			"&",
			[]any{"landed_cost_ok", "=", true},
			[]any{"type", "=", "service"},
		},
	},
	{
		"['&','&',('state','in',('sale','done')),('is_service','=',True),'|',('project_id','!=',False),('task_id','!=',False)]",
		[]any{
			"&",
			"&",
			[]any{"state", "in", []any{"sale", "done"}},
			[]any{"is_service", "=", true},
			"|",
			[]any{"project_id", "!=", false},
			[]any{"task_id", "!=", false},
		},
	},
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
