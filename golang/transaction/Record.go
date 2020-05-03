package transaction

import (
	. "github.com/saichler/gsql/golang/introspector"
	. "github.com/saichler/utils/golang"
	"reflect"
)

type Record struct {
	data map[string]reflect.Value
}

func (rec *Record) init() {
	if rec.data == nil {
		rec.data = make(map[string]reflect.Value)
	}
}

func (rec *Record) SetValue(key string, value reflect.Value) {
	rec.init()
	rec.data[key] = value
}

func (rec *Record) SetInterface(key string, any interface{}) {
	rec.init()
	rec.data[key] = reflect.ValueOf(any)
}

func (rec *Record) PrimaryIndex(pi *Index) string {
	result := NewStringBuilder("")
	for _, column := range pi.Columns() {
		val := rec.data[column.Name()]
		sv := ToString(val)
		result.Append(sv)
	}
	return result.String()
}

func (rec *Record) Data() map[string]reflect.Value {
	rec.init()
	return rec.data
}

func (rec *Record) Get(key string) reflect.Value {
	rec.init()
	return rec.data[key]
}
