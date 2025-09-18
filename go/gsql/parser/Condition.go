package parser

import (
	"bytes"
	"errors"
	"strings"

	"github.com/saichler/l8types/go/types/l8api"
)

type ConditionOperation string

const (
	And                 ConditionOperation = " and "
	Or                  ConditionOperation = " or "
	MAX_EXPRESSION_SIZE                    = 999999
)

func StringCondition(this *l8api.L8Condition) string {
	buff := &bytes.Buffer{}
	buff.WriteString("(")
	toString(this, buff)
	buff.WriteString(")")
	return buff.String()
}

func VisualizeCondition(this *l8api.L8Condition, lvl int) string {
	buff := &bytes.Buffer{}
	buff.WriteString(space(lvl))
	buff.WriteString("Condition\n")
	if this.Comparator != nil {
		buff.WriteString(VisualizeComparator(this.Comparator, lvl+1))
	}
	if this.Next != nil {
		buff.WriteString(space(lvl))
		buff.WriteString(strings.TrimSpace(this.Oper))
		buff.WriteString("\n")
		buff.WriteString(VisualizeCondition(this.Next, lvl))
	}
	return buff.String()
}

func toString(this *l8api.L8Condition, buff *bytes.Buffer) {
	if this.Comparator != nil {
		buff.WriteString(StringComparator(this.Comparator))
	}
	if this.Next != nil {
		buff.WriteString(this.Oper)
		toString(this.Next, buff)
	}
}

func NewCondition(ws string) (*l8api.L8Condition, error) {
	loc := MAX_EXPRESSION_SIZE
	var op ConditionOperation
	and := strings.Index(ws, string(And))
	if and != -1 {
		loc = and
		op = And
	}
	or := strings.Index(ws, string(Or))
	if or != -1 && or < loc {
		loc = or
		op = Or
	}

	condition := &l8api.L8Condition{}
	if loc == MAX_EXPRESSION_SIZE {
		cmpr, e := NewCompare(ws)
		if e != nil {
			return nil, e
		}
		condition.Comparator = cmpr
		return condition, nil
	}

	cmpr, e := NewCompare(ws[0:loc])
	if e != nil {
		return nil, e
	}

	condition.Comparator = cmpr
	condition.Oper = string(op)

	ws = ws[loc+len(op):]
	next, e := NewCondition(ws)
	if e != nil {
		return nil, e
	}

	condition.Next = next
	return condition, nil
}

func getLastConditionOp(ws string) (ConditionOperation, int, error) {
	loc := -1
	var op ConditionOperation

	and := strings.LastIndex(ws, string(And))
	if and > loc {
		op = And
		loc = and
	}

	or := strings.LastIndex(ws, string(Or))
	if or > loc {
		op = Or
		loc = or
	}

	if loc == -1 {
		return "", 0, errors.New("No last condition was found.")
	}
	return op, loc, nil
}

func getFirstConditionOp(ws string) (ConditionOperation, int, error) {
	loc := MAX_EXPRESSION_SIZE
	var op ConditionOperation
	and := strings.Index(ws, string(And))
	if and != -1 {
		loc = and
		op = And
	}
	or := strings.Index(ws, string(Or))
	if or != -1 && or < loc {
		loc = or
		op = Or
	}

	if loc == MAX_EXPRESSION_SIZE {
		return "", 0, errors.New("No first condition was found.")
	}

	return op, loc, nil
}
