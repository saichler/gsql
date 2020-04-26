package comparators

import (
	"reflect"
	"strconv"
	"strings"
)

type Equal struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewEqual() *Equal {
	c := &Equal{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = eqStringMatcher
	c.compares[reflect.Int] = eqIntMatcher
	c.compares[reflect.Int8] = eqIntMatcher
	c.compares[reflect.Int16] = eqIntMatcher
	c.compares[reflect.Int32] = eqIntMatcher
	c.compares[reflect.Int64] = eqIntMatcher
	c.compares[reflect.Uint] = eqUintMatcher
	c.compares[reflect.Uint8] = eqUintMatcher
	c.compares[reflect.Uint16] = eqUintMatcher
	c.compares[reflect.Uint32] = eqUintMatcher
	c.compares[reflect.Uint64] = eqUintMatcher
	c.compares[reflect.Ptr] = eqPtrMatcher
	return c
}

func (equal *Equal) Compare(left, right []reflect.Value) bool {
	return Compare(left, right, equal.compares, "Equal")
}

func Compare(left, right []reflect.Value, compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool, name string) bool {
	kind := getKind(left, right)
	compareFunc := compares[kind]
	if compareFunc == nil {
		panic("Cannot find compare func for:" + name + " Kind:" + kind.String())
	}
	for _, aside := range left {
		for _, zside := range right {
			if compareFunc(aside, zside) {
				return true
			}
		}
	}
	return false
}

func removeSingleQuote(value string) string {
	if strings.Contains(value, "'") {
		return value[1 : len(value)-1]
	}
	return value
}

func eqStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	if aside == "nil" && zside == "" {
		return true
	}
	if zside == "nil" && aside == "" {
		return true
	}
	splits := GetWildCardSubstrings(zside)
	if splits == nil {
		return aside == zside
	}
	for _, substr := range splits {
		if substr != "" && strings.Contains(aside, substr) {
			return true
		}
	}
	return false
}

func eqPtrMatcher(left, right reflect.Value) bool {
	if left.Kind() == reflect.Ptr && right.IsValid() && right.String() == "nil" && left.IsNil() {
		return true
	}
	if right.Kind() == reflect.Ptr && left.IsValid() && left.String() == "nil" && right.IsNil() {
		return true
	}
	return false
}

func eqIntMatcher(left, right reflect.Value) bool {
	aside, aok := getInt64(left)
	zside, zok := getInt64(right)

	if right.String() == "nil" && aok && aside == 0 {
		return true
	}

	if left.String() == "nil" && zok && zside == 0 {
		return true
	}

	if !aok || !zok {
		return false
	}

	return aside == zside
}

func eqUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside == zside
}

func getKind(aside, zside []reflect.Value) reflect.Kind {
	aSideKind := reflect.String
	zSideKind := reflect.String
	if len(aside) > 0 {
		aSideKind = aside[0].Kind()
	}
	if len(zside) > 0 {
		zSideKind = zside[0].Kind()
	}
	if aSideKind != reflect.String {
		return aSideKind
	} else if zSideKind != reflect.String {
		return zSideKind
	}
	return aSideKind
}

func getInt64(value reflect.Value) (int64, bool) {
	if value.Kind() != reflect.String {
		return value.Int(), true
	} else {
		i, e := strconv.Atoi(value.String())
		if e != nil {
			return 0, false
		}
		return int64(i), true
	}
}

func getUint64(value reflect.Value) (uint64, bool) {
	if value.Kind() != reflect.String {
		return value.Uint(), true
	} else {
		i, e := strconv.Atoi(value.String())
		if e != nil {
			return 0, false
		}
		return uint64(i), true
	}
}

func GetWildCardSubstrings(str string) []string {
	if !strings.Contains(str, "*") {
		return nil
	}
	return strings.Split(str, "*")
}
