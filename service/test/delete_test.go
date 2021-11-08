package test

import (
	"go-nosql-db/service"
	"go-nosql-db/util"
	"os"
	"testing"
)

func TestDeleARecords_searchBykeyAndValue(t *testing.T) {
	dataToSave := string(`{
        "home_address":
        {
            "line1": "wagholi",
            "state": "pune"
        },
        "subjets": ["hindi", "english", "math", "science"],
        "dob": "15-09-1986",
        "email": "ds.123@gmail.com",
        "lastName": "todelete",
        "name": "dharmendra"
    }`)
	os.Setenv("DB_DATA", "./")
	dbTable := "student_test"
	command := service.Persist{Table: dbTable, Data: dataToSave}
	command.Run()

	filter := service.Filter{Table: dbTable, FilterKey: "lastName", FilterValue: "todelete"}
	deleteCmd := service.Delete{Filter: filter}
	deleteCmd.Run()

	commandSelect := service.Select{Filter: filter, Projection: []string{"id"}}
	selectResult := commandSelect.Run()

	selectResultMinified, _ := util.MinifyJsonObject(selectResult)
	if string(selectResultMinified) != "null" {
		t.Errorf("record must Not exists")
	}
}

func TestDeleARecords_fromUnknownFile(t *testing.T) {
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test_unknown", FilterKey: "lastName", FilterValue: "todelete"}
	deleteCmd := service.Delete{Filter: filter}

	defer func() { recover() }()
	deleteCmd.Run()
	// Never reaches here if `Run()` panics.
	t.Errorf("did not panic")

}

func TestTruncateFileIfAllRecordsAreSelectedForDeletion(t *testing.T) {
	dataToSave := string(`{
        "home_address":
        {
            "line1": "wagholi",
            "state": "pune"
        },
        "subjets": ["hindi", "english", "math", "science"],
        "dob": "15-09-1986",
        "email": "ds.123@gmail.com",
        "lastName": "todelete",
        "name": "dharmendra"
    }`)
	os.Setenv("DB_DATA", "./")
	dbTable := "empty_file"
	jsonDb := util.JsonDb{Table: dbTable}
	jsonDb.Truncate()
	command := service.Persist{Table: dbTable, Data: dataToSave}
	command.Run()
	command.Run()

	filter := service.Filter{Table: dbTable, FilterKey: "name", FilterValue: "dharmendra"}
	deleteCmd := service.Delete{Filter: filter}
	deleteCmd.Run()

	commandSelect := service.Select{Filter: filter, Projection: []string{"id"}}
	selectResult := commandSelect.Run()

	selectResultMinified, _ := util.MinifyJsonObject(selectResult)
	fileSize := jsonDb.GetFileSize()
	if string(selectResultMinified) != "null" && fileSize > 0 {
		t.Errorf("record must Not exists")
	}
}
