package parser

import (
	"bytes"
	"errors"
	"strings"
)

func StringExpression(this *types.Expression) string {
	buff := bytes.Buffer{}
	if this.Condition != nil {
		buff.WriteString(StringCondition(this.Condition))
	} else {
		buff.WriteString("(")
	}
	if this.Child != nil {
		buff.WriteString(StringExpression(this.Child))
	}
	if this.Condition == nil {
		buff.WriteString(")")
	}
	if this.Next != nil {
		buff.WriteString(this.AndOr)
		buff.WriteString(StringExpression(this.Next))
	}
	return buff.String()
}

func VisualizeExpression(this *types.Expression, lvl int) string {
	buff := bytes.Buffer{}
	buff.WriteString(space(lvl))
	buff.WriteString("Expression\n")
	if this.Condition != nil {
		buff.WriteString(VisualizeCondition(this.Condition, lvl+1))
	}
	if this.Child != nil {
		buff.WriteString(VisualizeExpression(this.Child, lvl+1))
	}
	if this.Next != nil {
		buff.WriteString(space(lvl))
		buff.WriteString(strings.TrimSpace(this.AndOr))
		buff.WriteString("\n")
		buff.WriteString(VisualizeExpression(this.Next, lvl))
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

func parseExpression(ws string) (*types.Expression, error) {
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

func parseWithBrackets(ws string, bo int) (*types.Expression, error) {
	be, e := getBE(ws, bo)
	if e != nil {
		return nil, e
	}
	expr := &types.Expression{}
	child, e := parseExpression(ws[1:be])
	if e != nil {
		return nil, e
	}

	expr.Child = child

	if be < len(ws)-1 {
		op, loc, e := getFirstConditionOp(ws[be+1:])
		if e != nil {
			return nil, e
		}
		expr.AndOr = string(op)
		next, e := parseExpression(ws[be+1+loc+len(op):])
		if e != nil {
			return nil, e
		}
		expr.Next = next
	}
	return expr, nil
}

func parseBeforeBrackets(ws string, bo int) (*types.Expression, error) {
	prefix := ws[0:bo]
	op, loc, e := getLastConditionOp(prefix)
	if e != nil {
		return nil, e
	}
	expr, e := parseNoBrackets(prefix[0:loc])
	if e != nil {
		return nil, e
	}
	expr.AndOr = string(op)
	next, e := parseExpression(ws[bo:])
	if e != nil {
		return nil, e
	}
	expr.Next = next
	return expr, nil
}

func parseNoBrackets(ws string) (*types.Expression, error) {
	expr := &types.Expression{}
	condition, e := NewCondition(ws)
	if e != nil {
		return nil, e
	}
	expr.Condition = condition
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
