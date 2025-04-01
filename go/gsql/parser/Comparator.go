package parser

import (
	"bytes"
	"errors"
	"github.com/saichler/types/go/types"
	"strings"
)

type ComparatorOperation string

const (
	Eq    ComparatorOperation = "="
	Neq   ComparatorOperation = "!="
	GT    ComparatorOperation = ">"
	LT    ComparatorOperation = "<"
	GTEQ  ComparatorOperation = ">="
	LTEQ  ComparatorOperation = "<="
	IN    ComparatorOperation = " in "
	NOTIN ComparatorOperation = " not in "
)

var comparators = make([]ComparatorOperation, 0)

func initComparators() {
	if len(comparators) == 0 {
		comparators = append(comparators, GTEQ)
		comparators = append(comparators, LTEQ)
		comparators = append(comparators, Neq)
		comparators = append(comparators, Eq)
		comparators = append(comparators, GT)
		comparators = append(comparators, LT)
		comparators = append(comparators, NOTIN)
		comparators = append(comparators, IN)
	}
}

func StringComparator(this *types.Comparator) string {
	buff := bytes.Buffer{}
	buff.WriteString(this.Left)
	buff.WriteString(this.Oper)
	buff.WriteString(this.Right)
	return buff.String()
}

func VisualizeComparator(this *types.Comparator, lvl int) string {
	buff := bytes.Buffer{}
	buff.WriteString(space(lvl))
	buff.WriteString("Comparator (")
	buff.WriteString(this.Left)
	buff.WriteString(string(this.Oper))
	buff.WriteString(this.Right)
	buff.WriteString(")\n")
	return buff.String()
}

func NewCompare(ws string) (*types.Comparator, error) {
	for _, op := range comparators {
		loc := strings.Index(ws, string(op))
		if loc != -1 {
			cmp := &types.Comparator{}
			cmp.Left = strings.TrimSpace(strings.ToLower(ws[0:loc]))
			cmp.Right = strings.TrimSpace(strings.ToLower(ws[loc+len(op):]))
			cmp.Oper = string(op)
			if validateValue(cmp.Left) != "" {
				return nil, errors.New(validateValue(cmp.Left))
			}
			if validateValue(cmp.Right) != "" {
				return nil, errors.New(validateValue(cmp.Right))
			}
			return cmp, nil
		}
	}
	return nil, errors.New("Cannot find comparator operation in: " + ws)
}

func validateValue(ws string) string {
	bo := strings.Index(ws, "(")
	be := strings.Index(ws, ")")
	if bo != -1 || be != -1 {
		return "Value " + ws + " contain illegale brackets."
	}
	return ""
}
