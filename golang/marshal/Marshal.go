package marshal

import (
	. "github.com/saichler/gsql/golang/gschema"
	. "github.com/saichler/gsql/golang/introspector"
	. "github.com/saichler/gsql/golang/transaction"
	utils "github.com/saichler/utils/golang"
	"reflect"
	"strconv"
)

type Marshaler struct {
	introspector *Introspector
}

var marshalers = make(map[reflect.Kind]func(reflect.Value, *Introspector, *Transaction, Persistency, *RecordID) (reflect.Value, error))

func initMarshalers() {
	if len(marshalers) == 0 {
		marshalers[reflect.Ptr] = ptrMarshal
		marshalers[reflect.Struct] = structMarshal
		marshalers[reflect.Map] = mapMarshal
		marshalers[reflect.Slice] = sliceMarshal
		marshalers[reflect.String] = defaultMarshal
		marshalers[reflect.Int] = defaultMarshal
		marshalers[reflect.Int32] = defaultMarshal
		marshalers[reflect.Int64] = defaultMarshal
		marshalers[reflect.Uint] = defaultMarshal
		marshalers[reflect.Uint32] = defaultMarshal
		marshalers[reflect.Uint64] = defaultMarshal
		marshalers[reflect.Float64] = defaultMarshal
		marshalers[reflect.Float32] = defaultMarshal
		marshalers[reflect.Bool] = defaultMarshal
	}
}

func NewMarshaler(introspector *Introspector) *Marshaler {
	initMarshalers()
	m := &Marshaler{}
	m.introspector = introspector
	return m
}

func (m *Marshaler) Intospector() *Introspector {
	return m.introspector
}

func (m *Marshaler) Marshal(any interface{}, tx *Transaction, persistency Persistency) error {
	initMarshalers()
	if any == nil {
		return nil
	}
	value := reflect.ValueOf(any)
	value, err := marshal(value, m.introspector, tx, persistency, NewRecordID())
	return err
}

func marshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	marshalFunc := marshalers[value.Kind()]
	if marshalFunc == nil {
		panic("No Marshal Function for kind " + value.Kind().String())
	}
	return marshalFunc(value, introspector, tx, pr, recordID)
}

func ptrMarshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	if value.IsNil() {
		return reflect.ValueOf(""), nil
	}
	v := value.Elem()
	return marshal(v, introspector, tx, pr, recordID)
}

func structMarshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	tableName := value.Type().Name()
	//No need to do anything, nameless struct
	if tableName == "" {
		return value, nil
	}
	table := introspector.Table(tableName)
	if table == nil {
		panic("Table:" + tableName + " was not registered!")
	}

	rec := &Record{}
	rec.SetInterface(RECORD_LEVEL, recordID.Level())
	if table.Indexes().PrimaryIndex() == nil {
		rec.SetInterface(RECORD_ID, recordID.String())
		rec.SetInterface(RECORD_INDEX, recordID.Index)
	}
	subTables := make([]*Column, 0)
	for fieldName, column := range table.Columns() {
		if column.MetaData().ColumnTableName() == "" {
			fieldValue := value.FieldByName(fieldName)
			marshalValue, err := marshal(fieldValue, introspector, tx, pr, recordID)
			if err != nil {
				panic(err)
			}
			rec.SetValue(fieldName, marshalValue)
		} else {
			subTables = append(subTables, column)
		}
	}

	rid := ""

	if table.Indexes().PrimaryIndex() != nil {
		rid = rec.PrimaryIndex(table.Indexes().PrimaryIndex())
		tx.AddRecord(rec, tableName, rid)
	} else {
		tx.AddRecord(rec, tableName, recordID.String())
		rid = strconv.Itoa(recordID.Index)
	}

	for _, sbColumn := range subTables {
		isTable := sbColumn.MetaData().ColumnTableName() != ""
		if isTable {
			recordID.Add(table.Name(), sbColumn.Name(), rid)
		}
		fieldValue := value.FieldByName(sbColumn.Name())
		sbValue, err := marshal(fieldValue, introspector, tx, pr, recordID)
		if err != nil {
			return reflect.ValueOf(rec), err
		}
		rec.SetInterface(sbColumn.Name(), utils.ToString(sbValue))
		if isTable {
			recordID.Del()
		}
	}
	return reflect.ValueOf(recordID), nil
}

func sliceMarshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	if value.IsNil() {
		return reflect.ValueOf(""), nil
	}
	sb := utils.NewStringBuilder("[")
	for i := 0; i < value.Len(); i++ {
		recordID.Index = i
		v, e := marshal(value.Index(i), introspector, tx, pr, recordID)
		if e != nil {
			panic("Unable To marshal! " + e.Error())
		}
		if i != 0 {
			sb.Append(",")
		}
		sb.Append(utils.ToString(v))
	}
	sb.Append("]")
	return reflect.ValueOf(sb.String()), nil
}

func mapMarshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	if value.IsNil() {
		return reflect.ValueOf(""), nil
	}
	sb := utils.NewStringBuilder("[")
	mapKeys := value.MapKeys()
	for i, key := range mapKeys {
		mv := value.MapIndex(key)
		keyString := utils.ToString(key)
		recordID.Index = i
		v, e := marshal(mv, introspector, tx, pr, recordID)
		if e != nil {
			panic("Unable To marshal! " + e.Error())
		}
		if i > 0 {
			sb.Append(",")
		}
		sb.Append(keyString)
		sb.Append("=")
		sb.Append(utils.ToString(v))
	}
	sb.Append("]")
	return reflect.ValueOf(sb.String()), nil
}

func defaultMarshal(value reflect.Value, introspector *Introspector, tx *Transaction, pr Persistency, recordID *RecordID) (reflect.Value, error) {
	return value, nil
}
