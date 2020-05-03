package transaction

type Transaction struct {
	tableData map[string]map[string][]*Record
}

func (tx *Transaction) All() map[string]map[string][]*Record {
	return tx.tableData
}

func (tx *Transaction) AddRecord(record *Record, tableName, id string) {
	if tx.tableData == nil {
		tx.tableData = make(map[string]map[string][]*Record)
	}
	if tx.tableData[tableName] == nil {
		tx.tableData[tableName] = make(map[string][]*Record)
	}
	if tx.tableData[tableName][id] == nil {
		tx.tableData[tableName][id] = make([]*Record, 0)
	}
	tx.tableData[tableName][id] = append(tx.tableData[tableName][id], record)
}

func (tx *Transaction) Records(tableName string, id string) []*Record {
	if tx.tableData == nil {
		return nil
	}
	if tx.tableData[tableName] == nil {
		return nil
	}
	return tx.tableData[tableName][id]
}

func (tx *Transaction) AllRecords(tableName string) []*Record {
	if tx.tableData == nil {
		return nil
	}
	if tx.tableData[tableName] == nil {
		return nil
	}
	result := make([]*Record, 0)
	for _, recs := range tx.tableData[tableName] {
		for _, rec := range recs {
			result = append(result, rec)
		}
	}
	return result
}
