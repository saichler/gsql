package comparators

import (
	"reflect"
	"strings"
)

type LessThan struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewLessThan() *LessThan {
	c := &LessThan{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = ltStringMatcher
	c.compares[reflect.Int] = ltIntMatcher
	c.compares[reflect.Int8] = ltIntMatcher
	c.compares[reflect.Int16] = ltIntMatcher
	c.compares[reflect.Int32] = ltIntMatcher
	c.compares[reflect.Int64] = ltIntMatcher
	c.compares[reflect.Uint] = ltUintMatcher
	c.compares[reflect.Uint8] = ltUintMatcher
	c.compares[reflect.Uint16] = ltUintMatcher
	c.compares[reflect.Uint32] = ltUintMatcher
	c.compares[reflect.Uint64] = ltUintMatcher
	return c
}

func (lt *LessThan) Compare(left, right []reflect.Value) bool {
	return Compare(left, right, lt.compares, "Less Than")
}

func ltStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	return aside < zside
}

func ltIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside < zside
}

func ltUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside < zside
}
