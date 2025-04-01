package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/go/gsql/parser"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

type Condition struct {
	comparator *Comparator
	op         parser.ConditionOperation
	next       *Condition
}

func CreateCondition(c *types.Condition, rootTable *types.RNode, introspector common.IIntrospector) (*Condition, error) {
	condition := &Condition{}
	condition.op = parser.ConditionOperation(c.Oper)
	comp, e := CreateComparator(c.Comparator, rootTable, introspector)
	if e != nil {
		return nil, e
	}
	condition.comparator = comp
	if c.Next != nil {
		next, e := CreateCondition(c.Next, rootTable, introspector)
		if e != nil {
			return nil, e
		}
		condition.next = next
	}
	return condition, nil
}

func (condition *Condition) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("(")
	condition.toString(buff)
	buff.WriteString(")")
	return buff.String()
}

func (condition *Condition) toString(buff *bytes.Buffer) {
	if condition.comparator != nil {
		buff.WriteString(condition.comparator.String())
	}
	if condition.next != nil {
		buff.WriteString(string(condition.op))
		condition.next.toString(buff)
	}
}

func (condition *Condition) Match(root interface{}) (bool, error) {
	comp, e := condition.comparator.Match(root)
	if e != nil {
		return false, e
	}
	next := true
	if condition.op == parser.Or {
		next = false
	}
	if condition.next != nil {
		next, e = condition.next.Match(root)
		if e != nil {
			return false, e
		}
	}
	if condition.op == "" {
		return next && comp, nil
	}
	if condition.op == parser.And {
		return comp && next, nil
	}
	if condition.op == parser.Or {
		return comp || next, nil
	}
	return false, errors.New("Unsupported operation in match:" + string(condition.op))
}
