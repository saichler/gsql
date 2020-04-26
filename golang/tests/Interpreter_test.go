package tests

import (
	"fmt"
	. "github.com/saichler/gsql/golang/interpreter"
	. "github.com/saichler/gsql/golang/introspector"
	. "github.com/saichler/utils/golang"
	"github.com/saichler/utils/golang/tests"
	"reflect"
	"testing"
)

func TestQueryValidation(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	q, err := NewQuery(r, "Select String fRom nOde wHere (String=hello world or (String=hello orm and string2=myvalue and SlicePtrNoKey.string=192*))")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	fmt.Println(q.String())
}

func TestQueryMatch(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	_, err := NewQuery(r, "Select String fRom nOde wHere (String=hello world or (String=hello orm and String2=myvalue and SlicePtrNoKey.string=192*))")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
}

func TestFetchValue(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere (String=hello world or (String=hello orm and String2=myvalue and SlicePtrNoKey.string=192*))")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	columns := ormQuery.Columns()
	for _, c := range columns {
		val := c.ValueOf(reflect.ValueOf(node))
		if val[0].Kind() == reflect.String && val[0].String() != c.FieldName()+"-0" {
			t.Fail()
			Error("Expected value of " + c.FieldName() + " but found " + val[0].String())
			return
		}
	}
}

func TestMatchValue(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere (String=hello world or (String=hello orm and String2=myvalue and SlicePtrNoKey.string=192))")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]

	if ormQuery.Match(node) {
		t.Fail()
		Error("1) Expected not to match")
		return
	}

	node.String = "hello world"

	if !ormQuery.Match(node) {
		t.Fail()
		Error("2) Expected to match")
		return
	}

	node.String = "hello orm"
	node.String2 = "myvalue"
	node.Ptr.String = "193"

	if ormQuery.Match(node) {
		t.Fail()
		Error("3) Expected not to match")
		return
	}

	node.String = "hello orm"
	node.String2 = "myvalue"
	node.SlicePtrNoKey[0].String = "192"
	if !ormQuery.Match(node) {
		t.Fail()
		Error("4) Expected to match")
		return
	}

	ormQuery, err = NewQuery(r, "Select String fRom nOde wHere SlicePtrNoKey.string=192")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	if !ormQuery.Match(node) {
		t.Fail()
		Error("5) Expected to match")
		return
	}

	ormQuery, err = NewQuery(r, "Select String fRom nOde wHere SlicePtrNoKey.string=192 or SlicePtrNoKey.string=193")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	if !ormQuery.Match(node) {
		t.Fail()
		Error("6) Expected to match")
		return
	}
}

func TestMultiMatchValue(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere SlicePtrNoKey.string=192 or SlicePtrNoKey.string=192")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	if ormQuery.Match(node) {
		t.Fail()
		Error("Expected not to match")
		return
	}
}

func TestDeepMatchMultiValueMap(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere MapPrimary.String=Subnode6-0-index-0")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	node := tests.InitTestModel(1)[0]
	if !ormQuery.Match(node) {
		t.Fail()
		Error("Expected to match")
		return
	}
}

func TestDeepMatchMultiValueMap2(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere MapPrimary.String=Subnode6-0-index")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	node := tests.InitTestModel(1)[0]
	if ormQuery.Match(node) {
		t.Fail()
		Error("Expected not to match")
		return
	}
}

func TestDeepMatchMultiValueMap3(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere MapPrimary.String='MapPrimary'")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	node := tests.InitTestModel(1)[0]
	if ormQuery.Match(node) {
		t.Fail()
		Error("Expected not to match")
		return
	}
}

func TestDeepMatchMultiValueMap4(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere MapPrimary.String='Subnode6-0-index-0'")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}

	node := tests.InitTestModel(1)[0]
	if !ormQuery.Match(node) {
		t.Fail()
		Error("Expected to match")
		return
	}
}

func TestDeepMatchMultiValueSlice(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere SubNode2Slice.SliceInSlice.String=SubNode3-0-0-0")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	if !ormQuery.Match(node) {
		t.Fail()
		Error("Expected to match")
		return
	}
}

func TestDeepMatchMultiValueSlice2(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere SubNode2Slice.SliceInSlice.String=hello")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	if ormQuery.Match(node) {
		t.Fail()
		Error("Expected not to match")
		return
	}
}

func TestDeepMatchMultiValueSlice3(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere SubNode2Slice.SliceInSlice.Int=55")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	if ormQuery.Match(node) {
		t.Fail()
		Error("Expected not to match")
		return
	}
}

func TestDeepMatchMultiValueSlice4(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	ormQuery, err := NewQuery(r, "Select String fRom nOde wHere SubNode2Slice.SliceInSlice.Int=1")
	if err != nil {
		t.Fail()
		Error(err)
		return
	}
	node := tests.InitTestModel(1)[0]
	if !ormQuery.Match(node) {
		t.Fail()
		Error("Expected to match")
		return
	}
}

func TestSchemaTable(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	st, _ := r.GraphSchema().GraphSchemaNode("node")
	si := st.NewInstance(nil)
	newInstances := si.ValueOfOrCreate(reflect.ValueOf(nil))
	node := newInstances[0].AsValue().Interface().(*tests.Node)
	if node == nil {
		t.Fail()
		Error("Node is nil")
		return
	}

	st, _ = r.GraphSchema().GraphSchemaNode("node.Ptr")
	si = st.NewInstance(nil)
	newInstances = si.ValueOfOrCreate(reflect.ValueOf(nil))
	addr := newInstances[0].AsValue().Interface().(*tests.Node)
	if addr == nil {
		t.Fail()
		Error("NodeIp in node is nil")
		return
	}

	sf := r.GraphSchema().CreateAttribute("node.Ptr.string")
	expected := "my_expected_value"
	sf.SetValue(node, expected)
	if node.Ptr.String != expected {
		t.Fail()
		Error("Expected value")
		return
	}
}

func TestMapSetValueTableCreation(t *testing.T) {
	path := "node.MapPrimary[k1]"
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	schemaNode, keys := r.GraphSchema().GraphSchemaNode(path)
	if schemaNode == nil {
		t.Fail()
		Error("Could not find any schema node:" + path)
		return
	}
	si := schemaNode.NewInstance(keys)
	newInstances := si.ValueOfOrCreate(reflect.ValueOf(nil))
	node6 := newInstances[0].AsValue().Interface().(*tests.SubNode6)
	if node6 == nil {
		t.Fail()
		Error("Subnode 6 is nil")
		return
	}
}

func TestMapSetValue(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	schemaField := r.GraphSchema().CreateAttribute("node.MapPrimary[k1].String")

	if schemaField == nil {
		t.Fail()
		Error("Could not create schema field for node.MapPrimary[k1].String")
		return
	}

	expected := "My_Value"
	node := &tests.Node{}
	schemaField.SetValue(node, expected)

	if node.MapPrimary == nil {
		Error("map is nil")
		t.Fail()
		return
	}

	sb6 := node.MapPrimary["k1"]
	if sb6 == nil {
		Error("map entry is nil")
		t.Fail()
		return
	}

	if sb6.String != expected {
		Error("Expected String to be:" + expected)
		t.Fail()
		return
	}

	schemaField = r.GraphSchema().CreateAttribute("node.MapPrimary[k2].String")
	schemaField.SetValue(node, expected)

	sb6 = node.MapPrimary["k2"]
	if sb6 == nil {
		Error("SB6 2 map entry is nil")
		t.Fail()
		return
	}

	if sb6.String != expected {
		Error("Expected SB6 2 String to be:" + expected)
		t.Fail()
		return

	}
}

func TestSliceSetValue(t *testing.T) {
	r := NewIntrospector()
	r.Introspect(&tests.Node{}, nil)
	schemaField := r.GraphSchema().CreateAttribute("node.SubNode2Slice[1].SliceInSlice[1].String")

	if schemaField == nil {
		t.Fail()
		Error("Could not create schema field for node.SubNode2Slice[1].SliceInSlice[1].String")
		return
	}

	expected := "MyValue"
	node := &tests.Node{}
	schemaField.SetValue(node, expected)

	if node.SubNode2Slice == nil {
		Error("SubNode2Slice list is nil")
		t.Fail()
		return
	}

	if len(node.SubNode2Slice) < 2 {
		Error("Expected at least 2 SubNode2Slice instances.")
		t.Fail()
		return
	}

	ci := node.SubNode2Slice[1]
	if ci == nil {
		Error("SubNode2Slice 1 is nil")
		t.Fail()
		return
	}

	if ci.SliceInSlice == nil {
		Error("SliceInSlice list is nil")
		t.Fail()
		return
	}

	if len(ci.SliceInSlice) < 2 {
		Error("Expected SliceInSlice list to be at least 2")
		t.Fail()
		return
	}

	val := ci.SliceInSlice[1]

	if val == nil {
		Error("SliceInSlice 1 is nil")
		t.Fail()
		return
	}

	if val.String != expected {
		Error("Expected String to be:" + expected)
		t.Fail()
		return
	}

	schemaField = r.GraphSchema().CreateAttribute("node.SubNode2Slice[1].SliceInSlice[2].String")
	schemaField.SetValue(node, expected)

	if len(ci.SliceInSlice) < 3 {
		Error("Expected SliceInSlice list to be at least 3")
		t.Fail()
		return
	}

	val = ci.SliceInSlice[2]

	if val == nil {
		Error("Val 2 is nil")
		t.Fail()
		return
	}

	if val.String != expected {
		Error("Expected val 2 String to be:" + expected)
		t.Fail()
		return
	}
}
