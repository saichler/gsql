package comparators

import (
	"reflect"
	"strings"
)

type NotEqual struct {
	compares map[reflect.Kind]func(interface{}, interface{}) bool
}

func NewNotEqual() *NotEqual {
	c := &NotEqual{}
	c.compares = make(map[reflect.Kind]func(interface{}, interface{}) bool)
	c.compares[reflect.String] = noteqStringMatcher
	c.compares[reflect.Int] = noteqIntMatcher
	c.compares[reflect.Int8] = noteqIntMatcher
	c.compares[reflect.Int16] = noteqIntMatcher
	c.compares[reflect.Int32] = noteqIntMatcher
	c.compares[reflect.Int64] = noteqIntMatcher
	c.compares[reflect.Uint] = noteqUintMatcher
	c.compares[reflect.Uint8] = noteqUintMatcher
	c.compares[reflect.Uint16] = noteqUintMatcher
	c.compares[reflect.Uint32] = noteqUintMatcher
	c.compares[reflect.Uint64] = noteqUintMatcher
	return c
}

func (notequal *NotEqual) Compare(left, right interface{}) bool {
	return Compare(left, right, notequal.compares, "Not Equal")
}

func noteqStringMatcher(left, right interface{}) bool {
	aside := removeSingleQuote(strings.ToLower(left.(string)))
	zside := removeSingleQuote(strings.ToLower(right.(string)))
	return aside != zside
}

func noteqIntMatcher(left, right interface{}) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside != zside
}

func noteqUintMatcher(left, right interface{}) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside != zside
}
