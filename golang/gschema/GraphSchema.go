package gschema

import (
	"reflect"
	"strings"
	"sync"
)

type GraphSchema struct {
	graphSchemaNodes    map[string]*GraphSchemaNode
	fieldName2ProtoName map[string]map[string]string
	mutex               *sync.Mutex
}

type GraphSchemaProvider interface {
	Tables() []string
	GraphSchema() *GraphSchema
}

func NewGraphSchema() *GraphSchema {
	schema := &GraphSchema{}
	schema.graphSchemaNodes = make(map[string]*GraphSchemaNode)
	schema.fieldName2ProtoName = make(map[string]map[string]string)
	schema.mutex = &sync.Mutex{}
	return schema
}

func (graphSchema *GraphSchema) RegisterStruct(fieldName string, parent *GraphSchemaNode, structType reflect.Type) *GraphSchemaNode {
	if strings.Contains(fieldName, "XXX_") {
		//Skip Proto Buffer internal attributes
		return nil
	}
	graphSchemaNode := newGraphSchemaNode(fieldName, parent, structType)
	graphSchema.mutex.Lock()
	defer graphSchema.mutex.Unlock()
	old, ok := graphSchema.graphSchemaNodes[graphSchemaNode.ID()]
	if !ok {
		existingChildren := graphSchema.existingChildren(structType)
		graphSchema.graphSchemaNodes[graphSchemaNode.ID()] = graphSchemaNode
		for _, child := range existingChildren {
			newChild := newGraphSchemaNode(child.fieldName, graphSchemaNode, child.structType)
			graphSchema.graphSchemaNodes[newChild.ID()] = newChild
		}
		return graphSchemaNode
	} else {
		return old
	}
}

func (graphSchema *GraphSchema) existingChildren(structType reflect.Type) []*GraphSchemaNode {
	existingChildren := make([]*GraphSchemaNode, 0)
	for _, schemaNode := range graphSchema.graphSchemaNodes {
		found := false
		if schemaNode.Type().Name() == structType.Name() {
			found = true
			for _, childSchemaNode := range graphSchema.graphSchemaNodes {
				if childSchemaNode.parent != nil && childSchemaNode.parent.Type().Name() == structType.Name() {
					existingChildren = append(existingChildren, childSchemaNode)
				}
			}
		}
		if found {
			break
		}
	}
	return existingChildren
}

func (graphSchema *GraphSchema) GraphSchemaNodesMap() map[string]*GraphSchemaNode {
	cloneMap := make(map[string]*GraphSchemaNode)
	graphSchema.mutex.Lock()
	defer graphSchema.mutex.Unlock()
	for k, v := range graphSchema.graphSchemaNodes {
		cloneMap[k] = v
	}
	return cloneMap
}

func (graphSchema *GraphSchema) GraphSchemaNode(id string) (*GraphSchemaNode, *GraphKeys) {
	graphSchema.mutex.Lock()
	defer graphSchema.mutex.Unlock()
	snk, path := NewGraphKeys(id)
	return graphSchema.graphSchemaNodes[path], snk
}

func (graphSchema *GraphSchema) NewAttribute(fieldName string, instance *Instance) *Attribute {
	attribute := &Attribute{}
	attribute.fieldName = fieldName
	attribute.instance = instance
	return attribute
}

func (graphSchema *GraphSchema) CreateAttribute(id string) *Attribute {
	keys, path := NewGraphKeys(id)
	lastIndex := strings.LastIndex(path, ".")
	if lastIndex != -1 {
		tablePath := path[0:lastIndex]
		fieldName := path[lastIndex+1:]
		if fieldName == "" {
			return nil
		}
		graphSchemaNode, _ := graphSchema.GraphSchemaNode(tablePath)
		if graphSchemaNode == nil {
			return nil
		}
		for i := 0; i < graphSchemaNode.structType.NumField(); i++ {
			colName := graphSchemaNode.structType.Field(i).Name
			protoName := graphSchema.GetFieldProtoName(graphSchemaNode.structType.Name(), colName)
			if strings.ToLower(colName) == strings.ToLower(fieldName) || protoName == fieldName {
				return graphSchema.NewAttribute(colName, graphSchemaNode.NewInstance(keys))
			}
		}
		return nil
	}
	return nil
}

func (graphSchema *GraphSchema) AddFieldProtoName(structName, fieldName, protoName string) {
	graphSchema.mutex.Lock()
	defer graphSchema.mutex.Unlock()
	smap, ok := graphSchema.fieldName2ProtoName[structName]
	if !ok {
		smap = make(map[string]string)
		graphSchema.fieldName2ProtoName[structName] = smap
	}
	smap[fieldName] = protoName
}

func (graphSchema *GraphSchema) GetFieldProtoName(structName, fieldName string) string {
	graphSchema.mutex.Lock()
	defer graphSchema.mutex.Unlock()
	smap, ok := graphSchema.fieldName2ProtoName[structName]
	if !ok {
		return ""
	}
	return smap[fieldName]
}
