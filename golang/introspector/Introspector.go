package introspector

import (
	"github.com/saichler/gsql/golang/gschema"
	"reflect"
)

type Introspector struct {
	tables      map[string]*Table
	tablesList  []string
	annotations map[string]map[Annotation]*AnnotationEntry
	graphSchema *gschema.GraphSchema
}

func NewIntrospector() *Introspector {
	introspector := &Introspector{}
	introspector.annotations = make(map[string]map[Annotation]*AnnotationEntry)
	introspector.tables = make(map[string]*Table)
	introspector.tablesList = make([]string, 0)
	introspector.graphSchema = gschema.NewGraphSchema()
	return introspector
}

func (introspector *Introspector) Introspect(any interface{}, parent *gschema.GraphSchemaNode) {
	if introspector.graphSchema == nil {
		introspector.graphSchema = gschema.NewGraphSchema()
	}
	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.Slice {

	}
	schemaNode := introspector.graphSchema.RegisterStruct("", parent, value.Type())
	if schemaNode != nil {
		introspector.introspect(value.Type(), schemaNode)
	}
}

func (introspector *Introspector) introspect(structType reflect.Type, schemaNode *gschema.GraphSchemaNode) {
	table := introspector.Table(structType.Name())
	if table != nil {
		return
	}
	table = &Table{}
	table.structType = structType
	table.introspector = introspector
	introspector.tables[structType.Name()] = table
	table.inspect(schemaNode)
}

func (introspector *Introspector) Table(name string) *Table {
	if introspector.tables == nil {
		introspector.tables = make(map[string]*Table)
	}
	return introspector.tables[name]
}

func (introspector *Introspector) TablesMap() map[string]*Table {
	return introspector.tables
}

func (introspector *Introspector) Tables() []string {
	if introspector.tablesList == nil || len(introspector.tablesList) != len(introspector.tables) {
		introspector.tablesList = make([]string, 0)
		for tn, _ := range introspector.tables {
			introspector.tablesList = append(introspector.tablesList, tn)
		}
	}
	return introspector.tablesList
}

func (introspector *Introspector) GraphSchema() *gschema.GraphSchema {
	return introspector.graphSchema
}
