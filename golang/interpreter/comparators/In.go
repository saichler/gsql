package comparators

import (
	"reflect"
	"strconv"
	"strings"
)

type IN struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewIN() *IN {
	c := &IN{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = inStringMatcher
	c.compares[reflect.Int] = inIntMatcher
	c.compares[reflect.Int8] = inIntMatcher
	c.compares[reflect.Int16] = inIntMatcher
	c.compares[reflect.Int32] = inIntMatcher
	c.compares[reflect.Int64] = inIntMatcher
	c.compares[reflect.Uint] = inUintMatcher
	c.compares[reflect.Uint8] = inUintMatcher
	c.compares[reflect.Uint16] = inUintMatcher
	c.compares[reflect.Uint32] = inUintMatcher
	c.compares[reflect.Uint64] = inUintMatcher
	return c
}

func (in *IN) Compare(left, right []reflect.Value) bool {
	return Compare(left, right, in.compares, "In")
}

func inStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zsideList := strings.ToLower(right.String())
	values := getInStringList(zsideList)
	for _, v := range values {
		if aside == v {
			return true
		}
	}
	return false
}

func inIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}

	zsideList := strings.ToLower(right.String())

	values := getInStringList(zsideList)
	for _, v := range values {
		intV, e := strconv.Atoi(v)
		if e != nil {
			return false
		}
		if aside == int64(intV) {
			return true
		}
	}
	return false
}

func inUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}

	zsideList := strings.ToLower(right.String())

	values := getInStringList(zsideList)
	for _, v := range values {
		intV, e := strconv.Atoi(v)
		if e != nil {
			return false
		}
		if aside == uint64(intV) {
			return true
		}
	}
	return false
}

func getInStringList(str string) []string {
	index := strings.Index(str, "[")
	index2 := strings.Index(str, "]")
	lst := str[index+1 : index2]
	values := strings.Split(lst, ",")
	result := make([]string, 0)
	for _, v := range values {
		result = append(result, removeSingleQuote(v))
	}
	return result
}
