package interpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/gsql/golang/gschema"
	"github.com/saichler/gsql/golang/parser"
	"reflect"
)

type Expression struct {
	condition *Condition
	op        parser.ConditionOperation
	next      *Expression
	child     *Expression
}

func (expression *Expression) String() string {
	buff := bytes.Buffer{}
	if expression.condition != nil {
		buff.WriteString(expression.condition.String())
	} else {
		buff.WriteString("(")
	}
	if expression.child != nil {
		buff.WriteString(expression.child.String())
	}
	if expression.condition == nil {
		buff.WriteString(")")
	}
	if expression.next != nil {
		buff.WriteString(string(expression.op))
		buff.WriteString(expression.next.String())
	}
	return buff.String()
}

func CreateExpression(graphSchema *gschema.GraphSchema, mainTable *gschema.GraphSchemaNode, expr *parser.Expression) (*Expression, error) {
	if expr == nil {
		return nil, nil
	}
	ormExpr := &Expression{}
	ormExpr.op = expr.Operation()
	if expr.Condition() != nil {
		cond, e := CreateCondition(graphSchema, mainTable, expr.Condition())
		if e != nil {
			return nil, e
		}
		ormExpr.condition = cond
	}

	if expr.Child() != nil {
		child, e := CreateExpression(graphSchema, mainTable, expr.Child())
		if e != nil {
			return nil, e
		}
		ormExpr.child = child
	}

	if expr.Next() != nil {
		next, e := CreateExpression(graphSchema, mainTable, expr.Next())
		if e != nil {
			return nil, e
		}
		ormExpr.next = next
	}

	return ormExpr, nil
}

func (expression *Expression) Match(value reflect.Value) (bool, error) {
	cond := true
	child := true
	next := true
	var e error
	if expression.op == parser.Or {
		cond = false
		child = false
		next = false
	}
	if expression.condition != nil {
		cond, e = expression.condition.Match(value)
		if e != nil {
			return false, e
		}
	}
	if expression.child != nil {
		child, e = expression.child.Match(value)
		if e != nil {
			return false, e
		}
	}
	if expression.next != nil {
		next, e = expression.next.Match(value)
		if e != nil {
			return false, e
		}
	}
	if expression.op == "" {
		return child && next && cond, nil
	}
	if expression.op == parser.And {
		return child && next && cond, nil
	}
	if expression.op == parser.Or {
		return child || next || cond, nil
	}

	return false, errors.New("Unsupported operation in match:" + string(expression.op))
}
