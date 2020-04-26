package comparators

import (
	"reflect"
	"strings"
)

type GreaterThanOrEqual struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewGreaterThanOrEqual() *GreaterThanOrEqual {
	c := &GreaterThanOrEqual{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = gteqStringMatcher
	c.compares[reflect.Int] = gteqIntMatcher
	c.compares[reflect.Int8] = gteqIntMatcher
	c.compares[reflect.Int16] = gteqIntMatcher
	c.compares[reflect.Int32] = gteqIntMatcher
	c.compares[reflect.Int64] = gteqIntMatcher
	c.compares[reflect.Uint] = gteqUintMatcher
	c.compares[reflect.Uint8] = gteqUintMatcher
	c.compares[reflect.Uint16] = gteqUintMatcher
	c.compares[reflect.Uint32] = gteqUintMatcher
	c.compares[reflect.Uint64] = gteqUintMatcher
	return c
}

func (gteq *GreaterThanOrEqual) Compare(left, right []reflect.Value) bool {
	return Compare(left, right, gteq.compares, "Greater Than Or Equal")
}

func gteqStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	return aside >= zside
}

func gteqIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside >= zside
}

func gteqUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside >= zside
}
