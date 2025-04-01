package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/go/gsql/interpreter/comparators"
	"github.com/saichler/gsql/go/gsql/parser"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

type Comparator struct {
	left             string
	leftSchemaField  *properties.Property
	op               parser.ComparatorOperation
	right            string
	rightSchemaField *properties.Property
}

type Comparable interface {
	Compare(interface{}, interface{}) bool
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
		pid, _ := comparator.leftSchemaField.PropertyId()
		buff.WriteString(pid)
	} else {
		buff.WriteString(comparator.left)
	}
	buff.WriteString(string(comparator.op))
	if comparator.rightSchemaField != nil {
		pid, _ := comparator.rightSchemaField.PropertyId()
		buff.WriteString(pid)
	} else {
		buff.WriteString(comparator.right)
	}
	return buff.String()
}

func CreateComparator(c *types.Comparator, rootTable *types.RNode, introspector common.IIntrospector) (*Comparator, error) {
	initComparables()
	ormComp := &Comparator{}
	ormComp.op = parser.ComparatorOperation(c.Oper)
	ormComp.left = c.Left
	ormComp.right = c.Right
	leftProp := propertyPath(ormComp.left, rootTable.TypeName)
	rightProp := propertyPath(ormComp.right, rootTable.TypeName)
	ormComp.leftSchemaField, _ = properties.PropertyOf(leftProp, introspector)
	ormComp.rightSchemaField, _ = properties.PropertyOf(rightProp, introspector)

	if ormComp.leftSchemaField == nil && ormComp.rightSchemaField == nil {
		return nil, errors.New("No Field was found for comparator: " + c.String())
	}
	return ormComp, nil
}

func (comparator *Comparator) Match(root interface{}) (bool, error) {
	var leftValue interface{}
	var rightValue interface{}
	var err error
	if comparator.leftSchemaField != nil {
		leftValue, err = comparator.leftSchemaField.Get(root)
		if err != nil {
			return false, err
		}
	} else {
		leftValue = comparator.left
	}
	if comparator.rightSchemaField != nil {
		rightValue, err = comparator.rightSchemaField.Get(root)
		return false, err
	} else {
		rightValue = comparator.right
	}
	matcher := comparables[comparator.op]
	if matcher == nil {
		panic("No Matcher for: " + comparator.op + " operation.")
	}
	return matcher.Compare(leftValue, rightValue), nil
}
