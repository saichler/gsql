package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/golang/gschema"
	"github.com/saichler/gsql/golang/interpreter/comparators"
	"github.com/saichler/gsql/golang/parser"
	"reflect"
)

type Comparator struct {
	left             string
	leftSchemaField  *gschema.Attribute
	op               parser.ComparatorOperation
	right            string
	rightSchemaField *gschema.Attribute
}

type Comparable interface {
	Compare([]reflect.Value, []reflect.Value) bool
}

var comparables = make(map[parser.ComparatorOperation]Comparable)

func initComparables() {
	if len(comparables) == 0 {
		comparables[parser.Eq] = comparators.NewEqual()
		comparables[parser.Neq] = comparators.NewNotEqual()
		comparables[parser.NOTIN] = comparators.NewNotIN()
		comparables[parser.IN] = comparators.NewIN()
		comparables[parser.GT] = comparators.NewGreaterThan()
		comparables[parser.LT] = comparators.NewLessThan()
		comparables[parser.GTEQ] = comparators.NewGreaterThanOrEqual()
		comparables[parser.LTEQ] = comparators.NewLessThanOrEqual()
	}
}

func (comparator *Comparator) String() string {
	buff := bytes.Buffer{}
	if comparator.leftSchemaField != nil {
		buff.WriteString(comparator.leftSchemaField.ID())
	} else {
		buff.WriteString(comparator.left)
	}
	buff.WriteString(string(comparator.op))
	if comparator.rightSchemaField != nil {
		buff.WriteString(comparator.rightSchemaField.ID())
	} else {
		buff.WriteString(comparator.right)
	}
	return buff.String()
}

func CreateComparator(graphSchema *gschema.GraphSchema, mainTable *gschema.GraphSchemaNode, c *parser.Comparator) (*Comparator, error) {
	initComparables()
	ormComp := &Comparator{}
	ormComp.op = c.Operation()
	ormComp.left = c.Left()
	ormComp.right = c.Right()
	ormComp.leftSchemaField = graphSchema.CreateAttribute(mainTable.CreateFieldID(ormComp.left))
	ormComp.rightSchemaField = graphSchema.CreateAttribute(mainTable.CreateFieldID(ormComp.right))

	if ormComp.leftSchemaField == nil && ormComp.rightSchemaField == nil {
		return nil, errors.New("No Field was found for comparator:" + c.String())
	}
	return ormComp, nil
}

func (comparator *Comparator) Match(value reflect.Value) (bool, error) {
	var leftValue []reflect.Value
	var rightValue []reflect.Value
	if comparator.leftSchemaField != nil {
		leftValue = comparator.leftSchemaField.ValueOf(value)
	} else {
		leftValue = []reflect.Value{reflect.ValueOf(comparator.left)}
	}
	if comparator.rightSchemaField != nil {
		rightValue = comparator.rightSchemaField.ValueOf(value)
	} else {
		rightValue = []reflect.Value{reflect.ValueOf(comparator.right)}
	}
	matcher := comparables[comparator.op]
	if matcher == nil {
		panic("No Matcher for: " + comparator.op + " operation.")
	}
	return matcher.Compare(leftValue, rightValue), nil
}
