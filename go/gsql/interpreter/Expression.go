package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/go/gsql/parser"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

type Expression struct {
	condition *Condition
	operation parser.ConditionOperation
	next      *Expression
	child     *Expression
}

func (this *Expression) String() string {
	buff := bytes.Buffer{}
	if this.condition != nil {
		buff.WriteString(this.condition.String())
	} else {
		buff.WriteString("(")
	}
	if this.child != nil {
		buff.WriteString(this.child.String())
	}
	if this.condition == nil {
		buff.WriteString(")")
	}
	if this.next != nil {
		buff.WriteString(string(this.operation))
		buff.WriteString(this.next.String())
	}
	return buff.String()
}

func CreateExpression(expr *types.Expression, rootTable *types.RNode, introspector common.IIntrospector) (*Expression, error) {
	if expr == nil {
		return nil, nil
	}
	ormExpr := &Expression{}
	ormExpr.operation = parser.ConditionOperation(expr.AndOr)
	if expr.Condition != nil {
		cond, e := CreateCondition(expr.Condition, rootTable, introspector)
		if e != nil {
			return nil, e
		}
		ormExpr.condition = cond
	}

	if expr.Child != nil {
		child, e := CreateExpression(expr.Child, rootTable, introspector)
		if e != nil {
			return nil, e
		}
		ormExpr.child = child
	}

	if expr.Next != nil {
		next, e := CreateExpression(expr.Next, rootTable, introspector)
		if e != nil {
			return nil, e
		}
		ormExpr.next = next
	}

	return ormExpr, nil
}

func (this *Expression) Match(root interface{}) (bool, error) {
	cond := true
	child := true
	next := true
	var e error
	if this.operation == parser.Or {
		cond = false
		child = false
		next = false
	}
	if this.condition != nil {
		cond, e = this.condition.Match(root)
		if e != nil {
			return false, e
		}
	}
	if this.child != nil {
		child, e = this.child.Match(root)
		if e != nil {
			return false, e
		}
	}
	if this.next != nil {
		next, e = this.next.Match(root)
		if e != nil {
			return false, e
		}
	}
	if this.operation == "" {
		return child && next && cond, nil
	}
	if this.operation == parser.And {
		return child && next && cond, nil
	}
	if this.operation == parser.Or {
		return child || next || cond, nil
	}

	return false, errors.New("Unsupported operation in match:" + string(this.operation))
}

func (this *Expression) Condition() common.ICondition {
	return this.condition
}

func (this *Expression) Operator() string {
	return string(this.operation)
}

func (this *Expression) Next() common.IExpression {
	return this.next
}

func (this *Expression) Child() common.IExpression {
	return this.child
}
