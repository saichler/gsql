package interpreter

import "reflect"

func (query *Query) RequestedDataOnly(any interface{}) interface{} {
	if len(query.columns) == 0 {
		return any
	}
	value := reflect.ValueOf(any)
	result := query.newInstance(any)
	for _, attribute := range query.columns {
		attrValues := attribute.ValueOf(value)
		if attrValues != nil {
			for _, attrValue := range attrValues {
				attribute.SetValue(result, attrValue.Interface())
			}
		}
	}
	return result
}

func (query *Query) newInstance(any interface{}) interface{} {
	t := reflect.ValueOf(any).Elem().Type()
	return reflect.New(t).Interface()
}
