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
	operation  parser.ConditionOperation
	next       *Condition
}

func CreateCondition(c *types.Condition, rootTable *types.RNode, introspector common.IIntrospector) (*Condition, error) {
	condition := &Condition{}
	condition.operation = parser.ConditionOperation(c.Oper)
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

func (this *Condition) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("(")
	this.toString(buff)
	buff.WriteString(")")
	return buff.String()
}

func (this *Condition) toString(buff *bytes.Buffer) {
	if this.comparator != nil {
		buff.WriteString(this.comparator.String())
	}
	if this.next != nil {
		buff.WriteString(string(this.operation))
		this.next.toString(buff)
	}
}

func (this *Condition) Match(root interface{}) (bool, error) {
	comp, e := this.comparator.Match(root)
	if e != nil {
		return false, e
	}
	next := true
	if this.operation == parser.Or {
		next = false
	}
	if this.next != nil {
		next, e = this.next.Match(root)
		if e != nil {
			return false, e
		}
	}
	if this.operation == "" {
		return next && comp, nil
	}
	if this.operation == parser.And {
		return comp && next, nil
	}
	if this.operation == parser.Or {
		return comp || next, nil
	}
	return false, errors.New("Unsupported operation in match:" + string(this.operation))
}

func (this *Condition) Comparator() common.IComparator {
	return this.comparator
}

func (this *Condition) Operator() string {
	return string(this.operation)
}

func (this *Condition) Next() common.ICondition {
	return this.next
}
