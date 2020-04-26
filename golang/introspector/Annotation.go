package introspector

import (
	"bytes"
	"errors"
	"reflect"
)

type AnnotationEntry struct {
	annotation    Annotation
	structName    string
	attributeName string
	value1        string
	value2        string
}

func (introspector *Introspector) AddAnnotation(annotation Annotation, any interface{}, attributeName, value1, value2 string) error {
	if any == nil {
		return errors.New("Cannot add annotation for a nil struct")
	}
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	structName := v.Type().Name()
	_, ok := v.Type().FieldByName(attributeName)
	if !ok {
		return errors.New("Cannot find attribute " + attributeName + " inside struct " + structName)
	}
	ae := &AnnotationEntry{}
	ae.structName = structName
	ae.attributeName = attributeName
	ae.annotation = annotation
	ae.value1 = value1
	ae.value2 = value2

	key := structName + "." + attributeName
	if introspector.annotations[key] == nil {
		introspector.annotations[key] = make(map[Annotation]*AnnotationEntry)
	}
	introspector.annotations[key][annotation] = ae
	return nil
}

func (introspector *Introspector) AsTags(key string) string {
	buff := bytes.Buffer{}
	annotations, ok := introspector.annotations[key]
	if ok {
		for annotaion, data := range annotations {
			buff.WriteString(string(annotaion))
			buff.WriteString("=")
			buff.WriteString(data.value1)
			if data.value2 != "" {
				buff.WriteString(":")
				buff.WriteString(data.value2)
			}
			buff.WriteString(" ")
		}
	}
	return buff.String()
}
