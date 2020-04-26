package gschema

import (
	"bytes"
	"reflect"
	"strings"
)

type GraphSchemaNode struct {
	fieldName  string
	id         string
	parent     *GraphSchemaNode
	structType reflect.Type
}

func newGraphSchemaNode(fieldName string, parent *GraphSchemaNode, structType reflect.Type) *GraphSchemaNode {
	gsn := &GraphSchemaNode{}
	gsn.parent = parent
	gsn.fieldName = fieldName
	gsn.id = strings.ToLower(fieldName)
	gsn.structType = structType
	if gsn.parent == nil {
		gsn.id = strings.ToLower(gsn.structType.Name())
	}
	return gsn
}

func (graphSchemaNode *GraphSchemaNode) FieldName() string {
	return graphSchemaNode.fieldName
}

func (graphSchemaNode *GraphSchemaNode) Parent() *GraphSchemaNode {
	return graphSchemaNode.parent
}

func (graphSchemaNode *GraphSchemaNode) Type() reflect.Type {
	return graphSchemaNode.structType
}

func (graphSchemaNode *GraphSchemaNode) NewInterface() interface{} {
	return graphSchemaNode.NewValue().Interface()
}

func (graphSchemaNode *GraphSchemaNode) NewValue() reflect.Value {
	return reflect.New(graphSchemaNode.structType)
}

func (graphSchemaNode *GraphSchemaNode) ID() string {
	if graphSchemaNode.parent == nil {
		return graphSchemaNode.id
	}
	buff := bytes.Buffer{}
	buff.WriteString(graphSchemaNode.parent.ID())
	buff.WriteString(".")
	buff.WriteString(graphSchemaNode.id)
	return buff.String()
}

func (graphSchemaNode *GraphSchemaNode) NewInstance(graphKeys *GraphKeys) *Instance {
	if graphSchemaNode.parent == nil {
		instance := &Instance{}
		instance.graphSchemaNode = graphSchemaNode
		return instance
	}
	parent := graphSchemaNode.parent.NewInstance(graphKeys)
	instance := &Instance{}
	instance.graphSchemaNode = graphSchemaNode
	instance.parent = parent
	if graphKeys != nil {
		key := graphKeys.Key(instance.graphSchemaNode.ID())
		if key != "" {
			instance.key = reflect.ValueOf(key)
		}
	}
	return instance
}

func (graphSchemaNode *GraphSchemaNode) CreateFieldID(fieldName string) string {
	buff := bytes.Buffer{}
	buff.WriteString(graphSchemaNode.id)
	buff.WriteString(".")
	buff.WriteString(fieldName)
	return buff.String()
}
