package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/go/gsql/parser"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"reflect"
	"strings"
)

type Query struct {
	rootTable  *types.RNode
	columns    map[string]*properties.Property
	where      *Expression
	sortBy     string
	descending bool
	limit      int32
	page       int32
	matchCase  bool
	resources  common.IResources
	query      *types.Query
}

func NewFromQuery(query *types.Query, resources common.IResources) (*Query, error) {
	iQuery := &Query{}
	iQuery.columns = make(map[string]*properties.Property)
	iQuery.descending = query.Descending
	iQuery.matchCase = query.MatchCase
	iQuery.page = query.Page
	iQuery.limit = query.Limit
	iQuery.sortBy = query.SortBy
	iQuery.resources = resources
	iQuery.query = query

	err := iQuery.initTables(query)
	if err != nil {
		return nil, err
	}

	err = iQuery.initColumns(query, resources.Introspector())
	if err != nil {
		return nil, err
	}

	rootTable := iQuery.RootTable()
	if rootTable == nil {
		return nil, errors.New("root table is nil")
	}

	expr, err := CreateExpression(query.Criteria, rootTable, resources.Introspector())
	if err != nil {
		return nil, err
	}
	iQuery.where = expr

	return iQuery, nil
}

func NewQuery(gsql string, resources common.IResources) (*Query, error) {
	pQuery, err := parser.NewQuery(gsql, resources.Logger())
	if err != nil {
		return nil, err
	}
	return NewFromQuery(pQuery.Query(), resources)
}

func (this *Query) Query() *types.Query {
	return this.query
}

func (this *Query) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("Select ")
	first := true

	for _, column := range this.columns {
		if !first {
			buff.WriteString(", ")
		}
		id, _ := column.PropertyId()
		buff.WriteString(id)
		first = false
	}

	buff.WriteString(" From ")
	buff.WriteString(this.rootTable.TypeName)

	if this.where != nil {
		buff.WriteString(" Where ")
		buff.WriteString(this.where.String())
	}
	return buff.String()
}

func (this *Query) RootTable() *types.RNode {
	return this.rootTable
}

func (this *Query) Columns() map[string]*properties.Property {
	return this.columns
}

func (this *Query) OnlyTopLevel() bool {
	return true
}

func (this *Query) Descending() bool {
	return this.descending
}

func (this *Query) MatchCase() bool {
	return this.matchCase
}

func (this *Query) Page() int32 {
	return this.page
}

func (this *Query) Limit() int32 {
	return this.limit
}

func (this *Query) SortBy() string {
	return this.sortBy
}

func (this *Query) initTables(query *types.Query) error {
	node, ok := this.resources.Introspector().Node(query.RootType)
	if !ok {
		return this.resources.Logger().Error("Cannot find node for table ", query.RootType)
	}
	this.rootTable = node
	return nil
}

func (this *Query) initColumns(query *types.Query, introspector common.IIntrospector) error {
	if query.Properties != nil && len(query.Properties) == 1 && query.Properties[0] == "*" {
		return nil
	} else {
		for _, col := range query.Properties {
			propPath := propertyPath(col, this.rootTable.TypeName)
			prop, err := properties.PropertyOf(propPath, introspector)
			if err != nil {
				return this.resources.Logger().Error("cannot find property for col ", propPath, ":", err.Error())
			}
			this.columns[col] = prop
		}
	}
	return nil
}

func propertyPath(colName, rootTable string) string {
	rootTable = strings.ToLower(rootTable)
	if strings.Contains(colName, rootTable) {
		return colName
	}
	buff := bytes.Buffer{}
	buff.WriteString(rootTable)
	buff.WriteString(".")
	buff.WriteString(colName)
	return buff.String()
}

func (this *Query) match(root interface{}) (bool, error) {
	if root == nil {
		return false, nil
	}
	if this.rootTable == nil {
		return false, nil
	}
	if this.where == nil {
		return true, nil
	}
	return this.where.Match(root)
}

func (this *Query) Filter(list []interface{}, onlySelectedColumns bool) []interface{} {
	result := make([]interface{}, 0)
	for _, i := range list {
		if this.Match(i) {
			if !onlySelectedColumns || len(this.columns) == 0 {
				result = append(result, i)
			} else {
				result = append(result, this.cloneOnlyWithColumns(i))
			}
		}
	}
	return result
}

func (this *Query) Match(any interface{}) bool {
	m, e := this.match(any)
	if e != nil {
		this.resources.Logger().Error(e)
	}
	return m
}

func (this *Query) cloneOnlyWithColumns(any interface{}) interface{} {
	typ := reflect.ValueOf(any).Elem().Type()
	clone := reflect.New(typ).Interface()
	for _, column := range this.columns {
		v, _ := column.Get(any)
		column.Set(clone, v)
	}
	return clone
}

func (this *Query) CreateColumns(introspector common.IIntrospector) map[string]*properties.Property {
	result := make(map[string]*properties.Property)
	for attrName, attr := range this.rootTable.Attributes {
		if attr.IsStruct {
			continue
		}
		//@TODO - create a method to calc propertyid for node
		result[attrName], _ = properties.PropertyOf(attrName, introspector)
	}
	return result
}
