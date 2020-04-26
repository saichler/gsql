package parser

import (
	"bytes"
	"errors"
	"strings"
)

type ConditionOperation string

const (
	And                 ConditionOperation = " and "
	Or                  ConditionOperation = " or "
	MAX_EXPRESSION_SIZE                    = 999999
)

type Condition struct {
	comparator *Comparator
	op         ConditionOperation
	next       *Condition
}

func (condition *Condition) Comparator() *Comparator {
	return condition.comparator
}

func (condition *Condition) Operation() ConditionOperation {
	return condition.op
}

func (condition *Condition) Next() *Condition {
	return condition.next
}

func (condition *Condition) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("(")
	condition.toString(buff)
	buff.WriteString(")")
	return buff.String()
}

func (condition *Condition) Visualize(lvl int) string {
	buff := &bytes.Buffer{}
	buff.WriteString(space(lvl))
	buff.WriteString("Condition\n")
	if condition.comparator != nil {
		buff.WriteString(condition.comparator.Visualize(lvl + 1))
	}
	if condition.next != nil {
		buff.WriteString(space(lvl))
		buff.WriteString(strings.TrimSpace(string(condition.op)))
		buff.WriteString("\n")
		buff.WriteString(condition.next.Visualize(lvl))
	}
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

func NewCondition(ws string) (*Condition, error) {
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

	condition := &Condition{}
	if loc == MAX_EXPRESSION_SIZE {
		cmpr, e := NewCompare(ws)
		if e != nil {
			return nil, e
		}
		condition.comparator = cmpr
		return condition, nil
	}

	cmpr, e := NewCompare(ws[0:loc])
	if e != nil {
		return nil, e
	}

	condition.comparator = cmpr
	condition.op = op

	ws = ws[loc+len(op):]
	next, e := NewCondition(ws)
	if e != nil {
		return nil, e
	}

	condition.next = next
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
