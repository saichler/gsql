package gschema

import (
	"bytes"
	"reflect"
	"strings"
)

const (
	InstanceId = "_IID_"
)

//InstanceID identifies a sub instance in the model via a string
type InstanceID struct {
	name   string
	key    string
	isMap  bool
	parent *InstanceID
}

func NewInstanceID(name, key string, parent *InstanceID, isMap bool) *InstanceID {
	instanceId := &InstanceID{}
	instanceId.parent = parent
	instanceId.key = key
	instanceId.name = name
	instanceId.isMap = isMap
	return instanceId
}

func (instanceId *InstanceID) String() string {
	buff := &bytes.Buffer{}
	instanceId.string(buff)
	return buff.String()

}

func (instanceId *InstanceID) ParentKey() string {
	buff := &bytes.Buffer{}
	if instanceId.parent != nil {
		buff.WriteString(instanceId.parent.String())
		buff.WriteString(".")
		buff.WriteString(instanceId.name)
	}
	return buff.String()
}

func (instanceId *InstanceID) string(buff *bytes.Buffer) {
	if instanceId.parent != nil {
		instanceId.parent.string(buff)
		buff.WriteString(".")
	}
	buff.WriteString(instanceId.name)
	if instanceId.key != "" {
		buff.WriteString("[")
		if instanceId.isMap {
			buff.WriteString("{")
		}
		buff.WriteString(instanceId.key)
		if instanceId.isMap {
			buff.WriteString("}")
		}
		buff.WriteString("]")
	}
}

func (instanceId *InstanceID) FromString(str string) {
	lastIndex := strings.LastIndex(str, ".")
	var parent *InstanceID
	if lastIndex != -1 {
		parent = &InstanceID{}
		parent.FromString(str[0:lastIndex])
		str = str[lastIndex+1:]
	}

	instanceId.parent = parent
	index := strings.Index(str, "[")
	if index == -1 {
		instanceId.name = str
	} else {
		index2 := strings.Index(str, "]")
		instanceId.name = str[0:index]
		instanceId.key = str[index+1 : index2]
		if instanceId.key[0:1] == "{" {
			instanceId.isMap = true
			instanceId.key = instanceId.key[1 : len(instanceId.key)-1]
		}
	}
}

func (instanceId *InstanceID) AsValue() reflect.Value {
	return reflect.ValueOf(instanceId)
}
