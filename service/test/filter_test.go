package test

import (
	"go-nosql-db/service"
	"go-nosql-db/util"
	"os"
	"testing"
)

func TestSearchByFilterKeyAndValue(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "name", FilterValue: "pooja"}
	result, record := filter.Search()
	if len(result) == 0 || len(record.Positions) == 0 {
		t.Errorf("record must exists")
	}
}

func TestSearchByValue(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "", FilterValue: "pooja"}
	result, record := filter.Search()
	if len(result) == 0 || len(record.Positions) == 0 {
		t.Errorf("record must exists")
	}
}

func TestSearchByEmptyKeyAndEmptyValue(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "", FilterValue: ""}
	result, record := filter.Search()

	if len(result) > 0 || len(record.Positions) > 0 {
		t.Errorf("record must Not exists")
	}
}

func TestSearchByNestedKeyAndValue(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "address.line1", FilterValue: "shankargarh"}
	result, record := filter.Search()
	if len(result) == 0 || len(record.Positions) == 0 {
		t.Errorf("record must exists")
	}
}

func TestSearchWhenJsonFileIsEmpty(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	dbTable := "empty_file"
	jsonDb := util.JsonDb{Table: dbTable}
	jsonDb.Truncate()

	filter := service.Filter{Table: dbTable, FilterKey: "address.line1", FilterValue: "allahabad"}
	result, record := filter.Search()
	if len(result) > 0 || len(record.Positions) > 0 {
		t.Errorf("record must Not exists")
	}
}

func TestSearchDBTableNotFound(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test_notfound", FilterKey: "", FilterValue: ""}
	defer func() { recover() }()
	filter.Search()
	// Never reaches here if `Search()` panics.
	t.Errorf("did not panic")
}
