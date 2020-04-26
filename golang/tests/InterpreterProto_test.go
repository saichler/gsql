package tests

import (
	"github.com/saichler/gsql/golang/interpreter"
	"github.com/saichler/gsql/golang/introspector"
	utils "github.com/saichler/utils/golang"
	"github.com/saichler/utils/golang/tests"
	"testing"
)

func TestProtoQuery(t *testing.T) {
	pb := &tests.ProtoNode{}
	i := introspector.NewIntrospector()
	i.AddAnnotation(introspector.PrimaryKey, pb, "PString", "myindex", "0")
	i.Introspect(pb, nil)
	q, e := interpreter.NewQuery(i, "select subs.PString from ProtoNode")
	if e != nil {
		t.Fail()
		utils.Error(e)
		return
	}
	cols := q.Columns()
	for _, col := range cols {
		col.SetValue(pb, "Hello World")
	}
	if pb.Subs == nil {
		t.Fail()
		utils.Error("Expedcted subs not to be nil")
		return
	}
	if len(pb.Subs) == 0 {
		t.Fail()
		utils.Error("Expedcted len of subs to be greater than 0")
		return
	}
	if pb.Subs[0].PString != "Hello World" {
		t.Fail()
		utils.Error("Expedcted subs[0] to be Hello World")
		return
	}
}
