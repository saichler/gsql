package comparators

import (
	"reflect"
	"strings"
)

type GreaterThan struct {
	compares map[reflect.Kind]func(interface{}, interface{}) bool
}

func NewGreaterThan() *GreaterThan {
	c := &GreaterThan{}
	c.compares = make(map[reflect.Kind]func(interface{}, interface{}) bool)
	c.compares[reflect.String] = gtStringMatcher
	c.compares[reflect.Int] = gtIntMatcher
	c.compares[reflect.Int8] = gtIntMatcher
	c.compares[reflect.Int16] = gtIntMatcher
	c.compares[reflect.Int32] = gtIntMatcher
	c.compares[reflect.Int64] = gtIntMatcher
	c.compares[reflect.Uint] = gtUintMatcher
	c.compares[reflect.Uint8] = gtUintMatcher
	c.compares[reflect.Uint16] = gtUintMatcher
	c.compares[reflect.Uint32] = gtUintMatcher
	c.compares[reflect.Uint64] = gtUintMatcher
	return c
}

func (gt *GreaterThan) Compare(left, right interface{}) bool {
	return Compare(left, right, gt.compares, "Greater Than")
}

func gtStringMatcher(left, right interface{}) bool {
	aside := removeSingleQuote(strings.ToLower(left.(string)))
	zside := removeSingleQuote(strings.ToLower(right.(string)))
	return aside > zside
}

func gtIntMatcher(left, right interface{}) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside > zside
}

func gtUintMatcher(left, right interface{}) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside > zside
}
