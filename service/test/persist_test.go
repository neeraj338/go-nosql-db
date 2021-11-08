package test

import (
	"go-nosql-db/service"
	"go-nosql-db/util"
	"os"
	"testing"
)

func TestSaveARecords(t *testing.T) {
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
	command := service.Persist{Table: "student_test", Data: dataToSave}
	assignId := command.Run()
	var mockResultArr []interface{} = make([]interface{}, 0)
	mockResultArr = append(mockResultArr, assignId)
	mocksaveResult, _ := util.MinifyJsonObject(mockResultArr)

	filter := service.Filter{Table: "student_test", FilterKey: "name", FilterValue: "dharmendra"}
	commandSelect := service.Select{Filter: filter, Projection: []string{}}
	selectResult := commandSelect.Run()
	selectResultMinified, _ := util.MinifyJsonObject(selectResult)
	if string(mocksaveResult) != string(selectResultMinified) {
		t.Errorf("record must exists")
	}
}

func TestSaveARecords_To_NewTable(t *testing.T) {
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
	tableName := "empty_file"
	jsonDb := util.JsonDb{Table: tableName}
	jsonDb.Truncate()

	command := service.Persist{Table: tableName, Data: dataToSave}
	assignId := command.Run()
	var mockResultArr []interface{} = make([]interface{}, 0)
	mockResultArr = append(mockResultArr, assignId)
	mocksaveResult, _ := util.MinifyJsonObject(mockResultArr)

	filter := service.Filter{Table: tableName, FilterKey: "name", FilterValue: "dharmendra"}
	commandSelect := service.Select{Filter: filter, Projection: []string{}}
	selectResult := commandSelect.Run()
	selectResultMinified, _ := util.MinifyJsonObject(selectResult)
	if string(mocksaveResult) != string(selectResultMinified) {
		t.Errorf("record must exists")
	}
}

func TestSaveInvalidJson(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("DB_DATA", "./")

	invalidDataToSave := string(`{
        "home_address":
        {
            "line1": "wagholi",
            "state": "pune"
        }, 
    }`)

	os.Setenv("DB_DATA", "./")
	command := service.Persist{Table: "student_test", Data: invalidDataToSave}
	command.Run()
	// Never reaches here if `Run()` panics.
	t.Errorf("did not panic")
}
