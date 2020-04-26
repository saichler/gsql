package comparators

import (
	"reflect"
	"strconv"
	"strings"
)

type NotIN struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewNotIN() *NotIN {
	c := &NotIN{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = notinStringMatcher
	c.compares[reflect.Int] = notinIntMatcher
	c.compares[reflect.Int8] = notinIntMatcher
	c.compares[reflect.Int16] = notinIntMatcher
	c.compares[reflect.Int32] = notinIntMatcher
	c.compares[reflect.Int64] = notinIntMatcher
	c.compares[reflect.Uint] = notinUintMatcher
	c.compares[reflect.Uint8] = notinUintMatcher
	c.compares[reflect.Uint16] = notinUintMatcher
	c.compares[reflect.Uint32] = notinUintMatcher
	c.compares[reflect.Uint64] = notinUintMatcher
	return c
}

func (in *NotIN) Compare(left, right []reflect.Value) bool {
	return Compare(left, right, in.compares, "In")
}

func notinStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zsideList := strings.ToLower(right.String())
	values := getInStringList(zsideList)
	for _, v := range values {
		if aside == v {
			return false
		}
	}
	return true
}

func notinIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return true
	}

	zsideList := strings.ToLower(right.String())

	values := getInStringList(zsideList)
	for _, v := range values {
		intV, e := strconv.Atoi(v)
		if e != nil {
			return true
		}
		if aside == int64(intV) {
			return false
		}
	}
	return true
}

func notinUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return true
	}

	zsideList := strings.ToLower(right.String())

	values := getInStringList(zsideList)
	for _, v := range values {
		intV, e := strconv.Atoi(v)
		if e != nil {
			return true
		}
		if aside == uint64(intV) {
			return false
		}
	}
	return true
}
