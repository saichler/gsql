package gschema

import (
	"bytes"
	. "github.com/saichler/utils/golang"
	"reflect"
)

type Attribute struct {
	instance   *Instance
	fieldName  string
	MapKey     string
	SliceIndex int
}

func (attribute *Attribute) Instance() *Instance {
	return attribute.instance
}

func (attribute *Attribute) FieldName() string {
	return attribute.fieldName
}

func (attribute *Attribute) ID() string {
	buff := bytes.Buffer{}
	buff.WriteString(attribute.instance.ID())
	buff.WriteString(".")
	buff.WriteString(attribute.fieldName)
	return buff.String()
}

func (attribute *Attribute) ValueOf(root reflect.Value) []reflect.Value {
	structInstances := attribute.instance.ValueOf(root)
	results := make([]reflect.Value, 0)
	for _, instance := range structInstances {
		instanceValue := instance.value
		if instanceValue.Kind() == reflect.Ptr {
			if instanceValue.IsNil() {
				continue
			} else {
				instanceValue = instanceValue.Elem()
			}
		}

		if instanceValue.Kind() == reflect.Slice {
			for i := 0; i < instanceValue.Len(); i++ {
				elem := instanceValue.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				results = append(results, elem.FieldByName(attribute.fieldName))
			}
		} else {
			results = append(results, instanceValue.FieldByName(attribute.fieldName))
		}
	}
	return results
}

func (attribute *Attribute) SetValue(rootAny, any interface{}) {
	if attribute == nil {
		return
	}
	root := reflect.ValueOf(rootAny)
	structInstances := attribute.instance.valueOf(root, true)
	results := make([]reflect.Value, 0)
	for _, instance := range structInstances {
		instanceValue := instance.value
		if instanceValue.Kind() == reflect.Ptr {
			if instanceValue.IsNil() {
				continue
			} else {
				instanceValue = instanceValue.Elem()
			}
		}
		if instanceValue.Kind() == reflect.Slice {
			for i := 0; i < instanceValue.Len(); i++ {
				elem := instanceValue.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				results = append(results, elem.FieldByName(attribute.fieldName))
			}
		} else {
			v := instanceValue.FieldByName(attribute.fieldName)
			v1 := reflect.ValueOf(any)
			if v.Kind() == reflect.Map {
				if v.IsNil() {
					v.Set(reflect.MakeMap(v.Type()))
				}
				mapkey := reflect.ValueOf(attribute.MapKey)
				v.SetMapIndex(mapkey, v1)
			} else if v.Kind() == v1.Kind() {
				v.Set(v1)
			} else {
				Error("Value Type Does do not Match, expected:" + v.Kind().String() + " but got:" + v1.Kind().String() + " key=" + attribute.MapKey)
			}
		}
	}
}

func (attribute *Attribute) Get(root interface{}) []interface{} {
	if root == nil {
		return nil
	}
	result := make([]interface{}, 0)
	values := attribute.ValueOf(reflect.ValueOf(root))
	if values != nil {
		for _, v := range values {
			result = append(result, v.Interface())
		}
	}
	return result
}

func (attribute *Attribute) GetValue(root interface{}) interface{} {
	instances := attribute.Get(root)
	if len(instances) > 0 {
		return instances[0]
	}
	return nil
}

func (attribute *Attribute) GetString(root interface{}) string {
	instances := attribute.Get(root)
	if len(instances) > 0 {
		result, _ := instances[0].(string)
		return result
	}
	return ""
}

func (attribute *Attribute) GetInt(root interface{}) int {
	instances := attribute.Get(root)
	if len(instances) > 0 {
		result, _ := instances[0].(int)
		return result
	}
	return 0
}

func (attribute *Attribute) GetFloat(root interface{}) float64 {
	instances := attribute.Get(root)
	if len(instances) > 0 {
		result, _ := instances[0].(float64)
		return result
	}
	return 0
}

func (attribute *Attribute) Kind() reflect.Kind {
	parentType := attribute.instance.graphSchemaNode.structType
	field, ok := parentType.FieldByName(attribute.fieldName)
	if !ok {
		return reflect.Invalid
	}
	return field.Type.Kind()
}
