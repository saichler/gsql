package gschema

import (
	"bytes"
	"reflect"
	"strconv"
)

type Instance struct {
	graphSchemaNode *GraphSchemaNode
	key             reflect.Value
	value           reflect.Value
	parent          *Instance
}

func (instance *Instance) GraphSchemaNode() *GraphSchemaNode {
	return instance.graphSchemaNode
}

func (instance *Instance) Key() interface{} {
	return instance.key
}

func NewInstance(key, value reflect.Value, gsn *GraphSchemaNode, parent *Instance) *Instance {
	gi := &Instance{}
	gi.key = key
	gi.value = value
	gi.graphSchemaNode = gsn
	gi.parent = parent
	return gi
}

func (instance *Instance) newInstance(key, value reflect.Value) *Instance {
	return NewInstance(key, value, instance.graphSchemaNode, instance.parent)
}

func (instance *Instance) keyStringValue() string {
	if instance.key.IsValid() {
		kind := instance.key.Kind()
		if kind == reflect.Int || kind == reflect.Int64 || kind == reflect.Int32 {
			return strconv.Itoa(int(instance.key.Int()))
		} else if kind == reflect.Uint || kind == reflect.Uint64 || kind == reflect.Uint32 {
			return strconv.Itoa(int(instance.key.Uint()))
		}
		return instance.key.String()
	}
	return ""
}

func (instance *Instance) ID() string {
	key := instance.keyStringValue()
	if instance.parent == nil {
		if key == "" {
			return instance.graphSchemaNode.id
		} else {
			buff := bytes.Buffer{}
			buff.WriteString(instance.graphSchemaNode.id)
			buff.WriteString("[")
			buff.WriteString(key)
			buff.WriteString("]")
			return buff.String()
		}
	}
	buff := bytes.Buffer{}
	buff.WriteString(instance.parent.ID())
	buff.WriteString(".")
	buff.WriteString(instance.graphSchemaNode.id)
	if key != "" {
		buff.WriteString("[")
		buff.WriteString(key)
		buff.WriteString("]")
	}
	return buff.String()
}

func (instance *Instance) newMapValue(mapValue reflect.Value) []*Instance {
	results := make([]*Instance, 0)
	newItem := reflect.New(instance.graphSchemaNode.structType)
	newMap := reflect.MakeMap(reflect.MapOf(instance.key.Type(), newItem.Type()))
	newMap.SetMapIndex(instance.key, newItem)
	mapValue.Set(newMap)
	newInstance := NewInstance(instance.key, newItem, instance.graphSchemaNode, instance.parent)
	results = append(results, newInstance)
	return results
}

func (instance *Instance) mapValue(mapValue reflect.Value, createIfNil bool) []*Instance {
	if createIfNil && mapValue.IsNil() && instance.key.IsValid() {
		return instance.newMapValue(mapValue)
	}

	results := make([]*Instance, 0)

	if instance.key.IsValid() {
		mapItem := mapValue.MapIndex(instance.key)
		if !createIfNil && !mapItem.IsValid() {
			return results
		}
		if !mapItem.IsValid() {
			mapItem = reflect.New(instance.graphSchemaNode.structType)
			mapValue.SetMapIndex(instance.key, mapItem)
		}
		newInstance := NewInstance(instance.key, mapItem, instance.graphSchemaNode, instance.parent)
		results = append(results, newInstance)
		return results
	}

	mapKeys := mapValue.MapKeys()
	var newInstance *Instance
	for _, mapKey := range mapKeys {
		mapItem := mapValue.MapIndex(mapKey)
		if mapItem.Kind() == reflect.Ptr {
			if mapItem.IsNil() {
				continue
			} else {
				mapItem = mapItem.Elem()
			}
		}
		if newInstance == nil {
			newInstance = NewInstance(mapKey, mapItem, instance.graphSchemaNode, instance.parent)
		} else {
			newInstance = newInstance.newInstance(mapKey, mapItem)
		}
		results = append(results, newInstance)
	}
	return results
}

func (instance *Instance) newSliceValue(sliceValue reflect.Value) []*Instance {
	results := make([]*Instance, 0)
	newItem := reflect.New(instance.graphSchemaNode.structType)
	index, err := strconv.Atoi(instance.key.String())
	if err != nil {
		return results
	}
	newSlice := reflect.MakeSlice(reflect.SliceOf(newItem.Type()), index+1, index+1)
	newSlice.Index(index).Set(newItem)
	sliceValue.Set(newSlice)
	newInstance := NewInstance(instance.key, newItem, instance.graphSchemaNode, instance.parent)
	results = append(results, newInstance)
	return results
}

func (instance *Instance) enlargeSliceValue(sliceValue reflect.Value, newSize int) {
	newItem := reflect.New(instance.graphSchemaNode.structType)
	newSlice := reflect.MakeSlice(reflect.SliceOf(newItem.Type()), newSize, newSize)
	for i := 0; i < sliceValue.Len(); i++ {
		newSlice.Index(i).Set(sliceValue.Index(i))
	}
	sliceValue.Set(newSlice)
}

func (instance *Instance) sliceValue(sliceValue reflect.Value, createIfNil bool) []*Instance {
	if createIfNil && sliceValue.IsNil() && instance.key.IsValid() {
		return instance.newSliceValue(sliceValue)
	}

	results := make([]*Instance, 0)
	if !instance.key.IsValid() {
		instance.key = reflect.ValueOf("0")
	}

	if instance.key.IsValid() {
		index, err := strconv.Atoi(instance.key.String())
		if err != nil {
			return results
		}
		if sliceValue.Len() <= index {
			instance.enlargeSliceValue(sliceValue, index+1)
		}
		sliceItem := sliceValue.Index(index)
		if !createIfNil && (!sliceItem.IsValid() || sliceItem.IsNil()) {
			return results
		}
		if !sliceItem.IsValid() || sliceItem.IsNil() {
			sliceItem = reflect.New(instance.graphSchemaNode.structType)
			sliceValue.Index(index).Set(sliceItem)
		}
		newInstance := NewInstance(instance.key, sliceItem, instance.graphSchemaNode, instance.parent)
		results = append(results, newInstance)
		return results
	}

	var newInstance *Instance
	for i := 0; i < sliceValue.Len(); i++ {
		sliceItem := sliceValue.Index(i)
		if sliceItem.Kind() == reflect.Ptr {
			if sliceItem.IsNil() {
				continue
			} else {
				sliceItem = sliceItem.Elem()
			}
		}
		if newInstance == nil {
			newInstance = NewInstance(reflect.ValueOf(i), sliceItem, instance.graphSchemaNode, instance.parent)
		} else {
			newInstance = newInstance.newInstance(reflect.ValueOf(i), sliceItem)
		}
		results = append(results, newInstance)
	}
	return results
}

func (instance *Instance) valueOf(value reflect.Value, createIfNil bool) []*Instance {
	if instance.parent == nil {
		value := createIfNilOrInvalid(value, instance.graphSchemaNode.structType, createIfNil)
		newInstance := NewInstance(reflect.ValueOf(nil), value, instance.graphSchemaNode, nil)
		return []*Instance{newInstance}
	}
	parents := instance.parent.valueOf(value, createIfNil)
	results := make([]*Instance, 0)
	for _, parent := range parents {
		parentValue := parent.value
		if parentValue.Kind() == reflect.Ptr {
			parentValue = parentValue.Elem()
		}

		if parentValue.Kind() != reflect.Struct {
			return nil
		}

		myValue := parentValue.FieldByName(instance.graphSchemaNode.fieldName)
		if myValue.Kind() == reflect.Map {
			results = append(results, instance.mapValue(myValue, createIfNil)...)
		} else if myValue.Kind() == reflect.Slice {
			results = append(results, instance.sliceValue(myValue, createIfNil)...)
		} else {
			if myValue.IsNil() && createIfNil {
				v := createIfNilOrInvalid(myValue, instance.graphSchemaNode.structType, createIfNil)
				myValue.Set(v)
			}
			newInstance := instance.newInstance(reflect.ValueOf(nil), myValue)
			results = append(results, newInstance)
		}
	}
	return results
}

func (instance *Instance) ValueOf(value reflect.Value) []*Instance {
	return instance.valueOf(value, false)
}

func (instance *Instance) ValueOfOrCreate(value reflect.Value) []*Instance {
	return instance.valueOf(value, true)
}

func createIfNilOrInvalid(value reflect.Value, structType reflect.Type, createIfNil bool) reflect.Value {
	if createIfNil && (!value.IsValid() || value.Kind() == reflect.Ptr && value.IsNil()) {
		return reflect.New(structType)
	}
	return value
}

func (instance *Instance) AsValue() reflect.Value {
	return instance.value
}
