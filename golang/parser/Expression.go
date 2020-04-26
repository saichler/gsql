package parser

import (
	"bytes"
	"errors"
	"strings"
)

type Expression struct {
	condition *Condition
	op        ConditionOperation
	next      *Expression
	child     *Expression
}

func (expression *Expression) Condition() *Condition {
	return expression.condition
}

func (expression *Expression) Operation() ConditionOperation {
	return expression.op
}

func (expression *Expression) Next() *Expression {
	return expression.next
}

func (expression *Expression) Child() *Expression {
	return expression.child
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

func (expression *Expression) Visualize(lvl int) string {
	buff := bytes.Buffer{}
	buff.WriteString(space(lvl))
	buff.WriteString("Expression\n")
	if expression.condition != nil {
		buff.WriteString(expression.condition.Visualize(lvl + 1))
	}
	if expression.child != nil {
		buff.WriteString(expression.child.Visualize(lvl + 1))
	}
	if expression.next != nil {
		buff.WriteString(space(lvl))
		buff.WriteString(strings.TrimSpace(string(expression.op)))
		buff.WriteString("\n")
		buff.WriteString(expression.next.Visualize(lvl))
	}
	return buff.String()
}

func space(lvl int) string {
	buff := bytes.Buffer{}
	buff.WriteString("|")
	for i := 0; i < lvl; i++ {
		buff.WriteString("--")
	}
	return buff.String()
}

func parseExpression(ws string) (*Expression, error) {
	initComparators()
	ws = strings.TrimSpace(ws)
	bo := getBO(ws)
	if bo == -1 {
		return parseNoBrackets(ws)
	}

	if bo > 0 {
		return parseBeforeBrackets(ws, bo)
	}

	return parseWithBrackets(ws, bo)
}

func parseWithBrackets(ws string, bo int) (*Expression, error) {
	be, e := getBE(ws, bo)
	if e != nil {
		return nil, e
	}
	expr := &Expression{}
	child, e := parseExpression(ws[1:be])
	if e != nil {
		return nil, e
	}

	expr.child = child

	if be < len(ws)-1 {
		op, loc, e := getFirstConditionOp(ws[be+1:])
		if e != nil {
			return nil, e
		}
		expr.op = op
		next, e := parseExpression(ws[be+1+loc+len(op):])
		if e != nil {
			return nil, e
		}
		expr.next = next
	}
	return expr, nil
}

func parseBeforeBrackets(ws string, bo int) (*Expression, error) {
	prefix := ws[0:bo]
	op, loc, e := getLastConditionOp(prefix)
	if e != nil {
		return nil, e
	}
	expr, e := parseNoBrackets(prefix[0:loc])
	if e != nil {
		return nil, e
	}
	expr.op = op
	next, e := parseExpression(ws[bo:])
	if e != nil {
		return nil, e
	}
	expr.next = next
	return expr, nil
}

func parseNoBrackets(ws string) (*Expression, error) {
	expr := &Expression{}
	condition, e := NewCondition(ws)
	if e != nil {
		return nil, e
	}
	expr.condition = condition
	return expr, nil
}

func getBO(ws string) int {
	return strings.Index(ws, "(")
}

func getBE(ws string, bo int) (int, error) {
	count := 0
	for i := bo; i < len(ws); i++ {
		if byte(ws[i]) == byte('(') {
			count++
		} else if byte(ws[i]) == byte(')') {
			count--
		}
		if count == 0 {
			return i, nil
		}
	}
	return -1, errors.New("Missing close bracket in: " + ws)
}
