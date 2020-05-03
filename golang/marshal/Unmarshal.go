package marshal

import (
	. "github.com/saichler/gsql/golang/interpreter"
	. "github.com/saichler/gsql/golang/introspector"
	. "github.com/saichler/gsql/golang/transaction"
	. "github.com/saichler/habitat-orm/golang/common"
	"github.com/saichler/utils/golang"
	"reflect"
	"strconv"
	"strings"
)

var getters = make(map[reflect.Kind]func(*Column, *Record, *RecordID, *Transaction) reflect.Value)

func initGetters() {
	if len(getters) == 0 {
		getters[reflect.Ptr] = getPtr
		getters[reflect.String] = getDefault
		getters[reflect.Float32] = getDefault
		getters[reflect.Float64] = getDefault
		getters[reflect.Uint] = getDefault
		getters[reflect.Uint16] = getDefault
		getters[reflect.Uint32] = getDefault
		getters[reflect.Uint64] = getDefault
		getters[reflect.Int] = getDefault
		getters[reflect.Int16] = getDefault
		getters[reflect.Int32] = getDefault
		getters[reflect.Int64] = getDefault
		getters[reflect.Bool] = getDefault
		getters[reflect.Struct] = getStruct
		getters[reflect.Map] = getMap
		getters[reflect.Slice] = getSlice
	}
}

func (m *Marshaler) UnMarshal(query *Query, instrospector *Introspector, tx *Transaction) []interface{} {
	initGetters()
	instances := unmarshal(query, tx, instrospector, NewRecordID())
	result := make([]interface{}, len(instances))
	for i := 0; i < len(result); i++ {
		result[i] = instances[i].Interface()
	}
	return result
}

func unmarshal(query *Query, tx *Transaction, introspector *Introspector, id *RecordID) []reflect.Value {
	result := make([]reflect.Value, 0)
	mainTable, e := query.MainTable()
	if e != nil {
		return nil
	}
	table := introspector.Table(mainTable.Type().Name())
	records := tx.AllRecords(table.Name())
	for _, record := range records {
		if record.Get(RECORD_LEVEL).Int() == 0 || !query.OnlyTopLevel() {
			instance := table.NewInstance()
			result = append(result, instance)
			for _, column := range table.Columns() {
				field := instance.Elem().FieldByName(column.Name())
				key := record.PrimaryIndex(table.Indexes().PrimaryIndex())
				id.Add(table.Name(), column.Name(), key)
				fv := get(column, record, id, tx)
				instance.Elem()
				//fmt.Println(column.Name())
				if fv.IsValid() {
					field.Set(fv)
				}
				id.Del()
			}
		}
	}
	return result
}

func get(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	getter := getters[column.Type().Kind()]
	if getter == nil {
		panic("No Getter for kind:" + column.Type().String())
	}
	return getter(column, record, id, tx)
}

func getPtr(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	ptrKind := column.Type().Elem().Kind()
	if ptrKind == reflect.Struct {
		table := column.Table().Introspector().Table(column.Type().Elem().Name())
		if table.Indexes().PrimaryIndex() == nil {
			subRecords := tx.Records(column.MetaData().ColumnTableName(), id.String())
			if subRecords == nil {
				return reflect.ValueOf(nil)
			}
			return getStruct(column, subRecords[0], id, tx)
		}
		key := record.Get(column.Name()).String()
		if key == "" || key == "<invalid Value>" {
			return reflect.ValueOf(nil)
		}
		subRecords := tx.Records(column.MetaData().ColumnTableName(), key)
		return getStruct(column, subRecords[0], id, tx)
	} else if ptrKind == reflect.Slice {
		newSlice := reflect.MakeSlice(reflect.SliceOf(column.Type().Elem()), 0, 0)
		//@TODO implement
		return newSlice
	} else if ptrKind == reflect.Map {
		newMap := reflect.MakeMapWithSize(column.Type(), 0)
		//@TODO implement
		return newMap
	} else {
		panic("No Ptr Handle of:" + ptrKind.String())
	}
}

func getDefault(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	return record.Get(column.Name())
}

func getStruct(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	table := column.Table().Introspector().Table(column.MetaData().ColumnTableName())
	if table == nil {
		panic("Cannot find table name:" + column.MetaData().ColumnTableName())
	}
	instance := table.NewInstance()

	for _, c := range table.Columns() {
		fld := instance.Elem().FieldByName(c.Name())
		var key = "0"
		if table.Indexes().PrimaryIndex() != nil {
			key = record.PrimaryIndex(table.Indexes().PrimaryIndex())
		} else {
			value := record.Get(RECORD_ID)
			if !value.IsValid() {
				panic(table.Name() + ":" + column.Name())
			}
			if key == "" {

			}
		}
		id.Add(table.Name(), c.Name(), key)
		v := get(c, record, id, tx)
		if v.IsValid() {
			fld.Set(v)
		}
		id.Del()
	}
	return instance
}

func createMap2Index(tableName string, recordId *RecordID, tx *Transaction, dontHaveIndex bool) map[int64]*Record {
	if dontHaveIndex {
		map2Index := make(map[int64]*Record)
		id := recordId.String()
		recs := tx.Records(tableName, id)
		if recs == nil || len(recs) == 0 {
			return nil
		}
		for _, r := range recs {
			index := r.Get(RECORD_INDEX).Int()
			map2Index[index] = r
		}
		return map2Index
	}
	return nil
}

func getMap(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	value := record.Get(column.Name())
	vString := value.String()
	if value.IsValid() {
		if vString == "" {
			return reflect.ValueOf(nil)
		}
		if column.MetaData().ColumnTableName() == "" {
			return utils.FromString(vString, column.Type())
		} else {
			table := column.Table().Introspector().Table(column.MetaData().ColumnTableName())
			if table == nil {
				panic("No Table was found with name:" + column.MetaData().ColumnTableName())
			}
			elems := getElements(record, column)
			m := reflect.MakeMapWithSize(column.Type(), len(elems))
			dontHaveIndex := table.Indexes().PrimaryIndex() == nil
			map2Index := createMap2Index(table.Name(), id, tx, dontHaveIndex)

			for _, v := range elems {
				index := strings.Index(v, "=")
				key := v[0:index]
				val := v[index+1:]
				if dontHaveIndex {
					if map2Index == nil {
						return reflect.ValueOf(nil)
					}
					i, e := strconv.Atoi(val)
					if e != nil {
						panic("Index value in map is not int: " + e.Error())
					}
					r := map2Index[int64(i)]
					if r == nil {
						panic("Cannot find Map Index :" + val)
					}
					id.Index = i
					sval := getStruct(column, r, id, tx)
					m.SetMapIndex(utils.FromString(key, column.Type().Key()), sval)
				} else {
					recs := tx.Records(table.Name(), val)
					if recs == nil || len(recs) != 1 {
						panic("Cannot find records for key:" + val)
					} else {
						sval := getStruct(column, recs[0], NewRecordID(), tx)
						m.SetMapIndex(utils.FromString(key, column.Type().Key()), sval)
					}
				}
			}
			return m
		}
	}
	return reflect.ValueOf(nil)
}

func getSlice(column *Column, record *Record, id *RecordID, tx *Transaction) reflect.Value {
	value := record.Get(column.Name())
	vString := value.String()
	if value.IsValid() {
		if vString == "" {
			return reflect.ValueOf(nil)
		}
		if column.MetaData().ColumnTableName() == "" {
			return utils.FromString(vString, column.Type())
		} else {
			table := column.Table().Introspector().Table(column.MetaData().ColumnTableName())
			if table == nil {
				panic("No Table was found with name:" + column.MetaData().ColumnTableName())
			}
			if table.Indexes().PrimaryIndex() != nil {
				keys := getElements(record, column)
				newSlice := reflect.MakeSlice(column.Type(), len(keys), len(keys))
				for i, key := range keys {
					rec := tx.Records(table.Name(), key)[0]
					newSlice.Index(i).Set(getStruct(column, rec, NewRecordID(), tx))
				}
				return newSlice
			} else {
				recs := tx.Records(table.Name(), id.String())
				newSlice := reflect.MakeSlice(column.Type(), len(recs), len(recs))
				for i, rec := range recs {
					elem := getStruct(column, rec, id, tx)
					newSlice.Index(i).Set(elem)

				}
				return newSlice
			}
		}
	} else if column.MetaData().ColumnTableName() != "" {
		table := column.Table().Introspector().Table(column.MetaData().ColumnTableName())
		if table == nil {
			panic("No Table was found with name:" + column.MetaData().ColumnTableName())
		}
	}
	return reflect.ValueOf(nil)
}

func getElements(record *Record, column *Column) []string {
	elementsString := record.Get(column.Name()).String()
	elementsString = elementsString[1 : len(elementsString)-1]
	elems := strings.Split(elementsString, ",")
	return elems
}
