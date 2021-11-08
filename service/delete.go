package service

import (
	"go-nosql-db/util"
)

type Delete struct {
	Filter
}

func (d *Delete) Run() interface{} {
	dbTable := util.JsonDb{Table: d.Table}
	result, recordPositions := d.Search()

	if recordPositions.AllRecordsSelectedFromTable {
		dbTable.Truncate()
	} else {
		for _, r := range recordPositions.Positions {
			dbTable.DeleteRecord(r.Begin, r.Length, r.FilePosition)
		}
	}

	return result
}
