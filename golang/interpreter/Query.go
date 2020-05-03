package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/golang/gschema"
	"github.com/saichler/gsql/golang/parser"
	. "github.com/saichler/utils/golang"
	"reflect"
	"strings"
)

type Query struct {
	tables     map[string]*gschema.GraphSchemaNode
	columns    map[string]*gschema.Attribute
	where      *Expression
	sortBy     string
	descending bool
	limit      int
	page       int
	matchCase  bool
}

func (query *Query) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("Select ")
	first := true

	for _, column := range query.columns {
		if !first {
			buff.WriteString(", ")
		}
		buff.WriteString(column.ID())
		first = false
	}

	buff.WriteString(" From ")

	first = true
	for _, table := range query.tables {
		if !first {
			buff.WriteString(", ")
		}
		buff.WriteString(table.ID())
		first = false
	}

	if query.where != nil {
		buff.WriteString(" Where ")
		buff.WriteString(query.where.String())
	}
	return buff.String()
}

func (query *Query) Tables() map[string]*gschema.GraphSchemaNode {
	return query.tables
}

func (query *Query) Columns() map[string]*gschema.Attribute {
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

func (query *Query) initTables(provider gschema.GraphSchemaProvider, pq *parser.Query) error {
	for _, tableName := range pq.Tables() {
		found := false
		for _, name := range provider.Tables() {
			if strings.ToLower(name) == tableName {
				query.tables[tableName], _ = provider.GraphSchema().GraphSchemaNode(name)
				found = true
				break
			}
		}
		if !found {
			return errors.New("Could not find Struct " + tableName + " in Orm Registry.")
		}
	}
	return nil
}

func (query *Query) initColumns(provider gschema.GraphSchemaProvider, pq *parser.Query) error {
	mainTable, e := query.MainTable()
	if e != nil {
		return e
	}
	if pq.Columns() != nil && len(pq.Columns()) == 1 && pq.Columns()[0] == "*" {
		return nil
	} else {
		for _, col := range pq.Columns() {
			sf := provider.GraphSchema().CreateAttribute(mainTable.CreateFieldID(col))
			if sf == nil {
				return errors.New("Cannot find query field: " + col)
			}
			query.columns[col] = sf
		}
	}
	return nil
}

func NewQuery(provider gschema.GraphSchemaProvider, sql string) (*Query, error) {

	pQuery, err := parser.NewQuery(sql)
	if err != nil {
		return nil, err
	}
	iQuery := &Query{}
	iQuery.tables = make(map[string]*gschema.GraphSchemaNode)
	iQuery.columns = make(map[string]*gschema.Attribute)
	iQuery.descending = pQuery.Descending()
	iQuery.matchCase = pQuery.MatchCase()
	iQuery.page = pQuery.Page()
	iQuery.limit = pQuery.Limit()
	iQuery.sortBy = pQuery.SortBy()

	err = iQuery.initTables(provider, pQuery)
	if err != nil {
		return nil, err
	}

	err = iQuery.initColumns(provider, pQuery)
	if err != nil {
		return nil, err
	}

	mainTable, err := iQuery.MainTable()
	if err != nil {
		return nil, err
	}

	expr, err := CreateExpression(provider.GraphSchema(), mainTable, pQuery.Where())
	if err != nil {
		return nil, err
	}
	iQuery.where = expr

	return iQuery, nil
}

func (query *Query) MainTable() (*gschema.GraphSchemaNode, error) {
	for _, t := range query.tables {
		return t, nil
	}
	return nil, errors.New("No tables in query")
}

func (query *Query) match(value reflect.Value) (bool, error) {
	if !value.IsValid() {
		return false, nil
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return false, nil
		} else {
			value = value.Elem()
		}
	}
	tableName := strings.ToLower(value.Type().Name())
	table := query.tables[tableName]
	if table == nil {
		return false, nil
	}
	if query.where == nil {
		return true, nil
	}
	return query.where.Match(value)
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
	val := reflect.ValueOf(any)
	m, e := query.match(val)
	if e != nil {
		Error(e)
	}
	return m
}

func (query *Query) cloneOnlyWithColumns(any interface{}) interface{} {
	typ := reflect.ValueOf(any).Elem().Type()
	clone := reflect.New(typ).Interface()
	for _, column := range query.columns {
		v := column.GetValue(any)
		column.SetValue(clone, v)
	}
	return clone
}

func (query *Query) CreateColumns(provider gschema.GraphSchemaProvider) map[string]*gschema.Attribute {
	result := make(map[string]*gschema.Attribute)
	for _, tbl := range query.tables {
		for i := 0; i < tbl.Type().NumField(); i++ {
			fld := tbl.Type().Field(i)
			if fld.Type.Kind() != reflect.Slice && fld.Type.Kind() != reflect.Map && fld.Type.Kind() != reflect.Ptr {
				sf := provider.GraphSchema().CreateAttribute(tbl.CreateFieldID(fld.Name))
				if sf != nil {
					result[fld.Name] = sf
				}
			}
		}
	}
	return result
}
