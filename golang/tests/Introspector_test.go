package tests

import (
	"fmt"
	"github.com/saichler/gsql/golang/introspector"
	"github.com/saichler/utils/golang/tests"
	"testing"
)

func TestIntrospection(t *testing.T) {
	node := &tests.Node{}
	i := introspector.NewIntrospector()
	i.Introspect(node, nil)
	for _, gsn := range i.GraphSchema().GraphSchemaNodesMap() {
		fmt.Println(gsn.ID())
	}
}

func TestIntrospectionOnProto(t *testing.T) {
	node := &tests.ProtoNode{}
	i := introspector.NewIntrospector()
	i.AddAnnotation(introspector.PrimaryKey, node, "PString", "myindex", "0")
	i.Introspect(node, nil)
	for _, gsn := range i.GraphSchema().GraphSchemaNodesMap() {
		fmt.Println(gsn.ID())
	}
}
