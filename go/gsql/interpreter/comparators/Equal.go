package comparators

import (
	"reflect"
	"strconv"
	"strings"
)

type Equal struct {
	compares map[reflect.Kind]func(interface{}, interface{}) bool
}

func NewEqual() *Equal {
	c := &Equal{}
	c.compares = make(map[reflect.Kind]func(interface{}, interface{}) bool)
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

func (equal *Equal) Compare(left, right interface{}) bool {
	return Compare(left, right, equal.compares, "Equal")
}

func Compare(left, right interface{}, compares map[reflect.Kind]func(interface{}, interface{}) bool, name string) bool {
	kind := getKind(left, right)
	compareFunc := compares[kind]
	if compareFunc == nil {
		panic("Cannot find compare func for:" + name + " Kind:" + kind.String())
	}
	return compareFunc(left, right)
}

func removeSingleQuote(value string) string {
	if strings.Contains(value, "'") {
		return value[1 : len(value)-1]
	}
	return value
}

func eqStringMatcher(left, right interface{}) bool {
	vLeft := reflect.ValueOf(left)
	if vLeft.Kind() == reflect.Slice {
		for i := 0; i < vLeft.Len(); i++ {
			if eqStringMatcher(vLeft.Index(i).Interface(), right) {
				return true
			}
		}
		return false
	}
	aside := removeSingleQuote(strings.ToLower(left.(string)))
	zside := removeSingleQuote(strings.ToLower(right.(string)))
	if aside == "nil" && zside == "" {
		return true
	}
	if zside == "nil" && aside == "" {
		return true
	}
	if aside == "*" || zside == "*" {
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

func eqPtrMatcher(left, right interface{}) bool {
	if left == nil && right.(string) == "nil" {
		return true
	}
	if right == nil && left.(string) == "nil" {
		return true
	}
	return false
}

func eqIntMatcher(left, right interface{}) bool {
	aside, aok := getInt64(left)
	zside, zok := getInt64(right)

	rightValue, ok := right.(string)
	if ok && rightValue == "nil" && aok && aside == 0 {
		return true
	}

	leftValue, ok := left.(string)
	if ok && leftValue == "nil" && zok && zside == 0 {
		return true
	}

	if !aok || !zok {
		return false
	}

	return aside == zside
}

func eqUintMatcher(left, right interface{}) bool {
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

func getKind(aside, zside interface{}) reflect.Kind {
	aSideKind := reflect.String
	zSideKind := reflect.String

	asideValue := reflect.ValueOf(aside)
	zsideValue := reflect.ValueOf(zside)

	if asideValue.IsValid() {
		if asideValue.Kind() == reflect.Slice {
			if asideValue.Len() > 0 {
				aSideKind = asideValue.Index(0).Kind()
				if aSideKind == reflect.Interface {
					aSideKind = reflect.ValueOf(asideValue.Index(0).Interface()).Kind()
				}
			}
		} else {
			aSideKind = asideValue.Kind()
		}
	}

	if zsideValue.IsValid() {
		if zsideValue.Kind() == reflect.Slice {
			if zsideValue.Len() > 0 {
				zSideKind = zsideValue.Index(0).Kind()
				if zSideKind == reflect.Interface {
					zSideKind = reflect.ValueOf(zsideValue.Index(0).Interface()).Kind()
				}
			}
		} else {
			zSideKind = zsideValue.Kind()
		}
	}

	if aSideKind != reflect.String {
		return aSideKind
	} else if zSideKind != reflect.String {
		return zSideKind
	}
	return aSideKind
}

func getInt64(v interface{}) (int64, bool) {
	value := reflect.ValueOf(v)
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

func getUint64(v interface{}) (uint64, bool) {
	value := reflect.ValueOf(v)
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
