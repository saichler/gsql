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
	limit      int
	page       int
	matchCase  bool
	resources  common.IResources
}

func NewQuery(gsql string, resources common.IResources) (*Query, error) {
	pQuery, err := parser.NewQuery(gsql, resources.Logger())
	if err != nil {
		return nil, err
	}
	iQuery := &Query{}
	iQuery.columns = make(map[string]*properties.Property)
	iQuery.descending = pQuery.Descending()
	iQuery.matchCase = pQuery.MatchCase()
	iQuery.page = pQuery.Page()
	iQuery.limit = pQuery.Limit()
	iQuery.sortBy = pQuery.SortBy()
	iQuery.resources = resources

	err = iQuery.initTables(pQuery)
	if err != nil {
		return nil, err
	}

	err = iQuery.initColumns(pQuery, resources.Introspector())
	if err != nil {
		return nil, err
	}

	rootTable := iQuery.RootTable()
	if rootTable == nil {
		return nil, errors.New("root table is nil")
	}

	expr, err := CreateExpression(pQuery.Where(), rootTable, resources.Introspector())
	if err != nil {
		return nil, err
	}
	iQuery.where = expr

	return iQuery, nil
}

func (query *Query) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("Select ")
	first := true

	for _, column := range query.columns {
		if !first {
			buff.WriteString(", ")
		}
		id, _ := column.PropertyId()
		buff.WriteString(id)
		first = false
	}

	buff.WriteString(" From ")
	buff.WriteString(query.rootTable.TypeName)

	if query.where != nil {
		buff.WriteString(" Where ")
		buff.WriteString(query.where.String())
	}
	return buff.String()
}

func (query *Query) RootTable() *types.RNode {
	return query.rootTable
}

func (query *Query) Columns() map[string]*properties.Property {
	return query.columns
}

func (query *Query) OnlyTopLevel() bool {
	return true
}

func (query *Query) Descending() bool {
	return query.descending
}

func (query *Query) MatchCase() bool {
	return query.matchCase
}

func (query *Query) Page() int {
	return query.page
}

func (query *Query) Limit() int {
	return query.limit
}

func (query *Query) SortBy() string {
	return query.sortBy
}

func (query *Query) initTables(pq *parser.Query) error {
	node, ok := query.resources.Introspector().Node(pq.RootTable())
	if !ok {
		return query.resources.Logger().Error("Cannot find node for table ", pq.RootTable())
	}
	query.rootTable = node
	return nil
}

func (query *Query) initColumns(pq *parser.Query, introspector common.IIntrospector) error {
	if pq.Columns() != nil && len(pq.Columns()) == 1 && pq.Columns()[0] == "*" {
		return nil
	} else {
		for _, col := range pq.Columns() {
			propPath := propertyPath(col, query.rootTable.TypeName)
			prop, err := properties.PropertyOf(propPath, introspector)
			if err != nil {
				return query.resources.Logger().Error("cannot find property for col ", propPath, ":", err.Error())
			}
			query.columns[col] = prop
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

func (query *Query) match(root interface{}) (bool, error) {
	if root == nil {
		return false, nil
	}
	if query.rootTable == nil {
		return false, nil
	}
	if query.where == nil {
		return true, nil
	}
	return query.where.Match(root)
}

func (query *Query) Filter(list []interface{}, onlySelectedColumns bool) []interface{} {
	result := make([]interface{}, 0)
	for _, i := range list {
		if query.Match(i) {
			if !onlySelectedColumns || len(query.columns) == 0 {
				result = append(result, i)
			} else {
				result = append(result, query.cloneOnlyWithColumns(i))
			}
		}
	}
	return result
}

func (query *Query) Match(any interface{}) bool {
	m, e := query.match(any)
	if e != nil {
		query.resources.Logger().Error(e)
	}
	return m
}

func (query *Query) cloneOnlyWithColumns(any interface{}) interface{} {
	typ := reflect.ValueOf(any).Elem().Type()
	clone := reflect.New(typ).Interface()
	for _, column := range query.columns {
		v, _ := column.Get(any)
		column.Set(clone, v)
	}
	return clone
}

func (query *Query) CreateColumns(introspector common.IIntrospector) map[string]*properties.Property {
	result := make(map[string]*properties.Property)
	for attrName, attr := range query.rootTable.Attributes {
		if attr.IsStruct {
			continue
		}
		//@TODO - create a method to calc propertyid for node
		result[attrName], _ = properties.PropertyOf(attrName, introspector)
	}
	return result
}
