package gschema

import (
	"github.com/saichler/utils/golang"
)

const (
	RECORD_ID    = "_RID"
	RECORD_INDEX = "_RINDEX"
	RECORD_LEVEL = "_RLVL"
	NO_INDEX     = -9
)

type RecordID struct {
	entries  []*RecordIDEntry
	location int
	Index    int
}

type RecordIDEntry struct {
	tableName  string
	columnName string
	parentKey  string
}

func (rid *RecordIDEntry) String() string {
	result := utils.NewStringBuilder("[")
	result.Append(rid.tableName).Append(".")
	result.Append(rid.columnName).Append("=")
	result.Append(rid.parentKey)
	result.Append("]")
	return result.String()
}

func NewRecordID() *RecordID {
	rid := &RecordID{}
	rid.entries = make([]*RecordIDEntry, 0)
	rid.location = -1
	return rid
}

func (rid *RecordID) Add(tableName, columnName, parentKey string) {
	if rid.entries == nil {
		panic("RecordID was not created with NewRecordID method.")
	}
	ride := &RecordIDEntry{}
	ride.tableName = tableName
	ride.columnName = columnName
	ride.parentKey = parentKey
	rid.entries = append(rid.entries, ride)
	rid.location++
	rid.Index = NO_INDEX
}

func (rid *RecordID) SetParentKey(parentKey string) {
	rid.entries[rid.location].parentKey = parentKey
}

func (rid *RecordID) Del() {
	rid.entries = rid.entries[0:rid.location]
	rid.location--
}

func (rid *RecordID) String() string {
	sb := utils.NewStringBuilder("")
	for _, s := range rid.entries {
		sb.Append(s.String())
	}
	return sb.String()
}

func (rid *RecordID) Level() int {
	return len(rid.entries)
}
