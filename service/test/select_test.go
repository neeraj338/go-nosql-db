package test

import (
	"go-nosql-db/service"
	"go-nosql-db/util"
	"os"
	"testing"
)

func TestSelectByProjection(t *testing.T) {
	mockData, _ := util.ToJsonObject(string(`[{"a5bdf0df-432e-4272-69ea-87b13ba06a8a":{
			"address": {
				"line1": "shankargarh",
				"state": "uttar-pradesh"
			},
			"name": "pooja"
		}
	  }]`))
	minifyMock, _ := util.MinifyJsonObject(mockData)
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "name", FilterValue: "pooja"}
	command := service.Select{Filter: filter, Projection: []string{"name", "id", "address"}}
	result := command.Run()
	minifyResult, _ := util.MinifyJsonObject(result)

	if result == nil || string(minifyMock) != string(minifyResult) {
		t.Errorf("record must exists")
	}
}

func TestSelectByProjection_NestedJsonKeys(t *testing.T) {
	mockData, _ := util.ToJsonObject(string(`[{"a5bdf0df-432e-4272-69ea-87b13ba06a8a":{
			"address_line1": "shankargarh",
			"name": "pooja"
		}
	  }]`))
	minifyMock, _ := util.MinifyJsonObject(mockData)
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "name", FilterValue: "pooja"}
	command := service.Select{Filter: filter, Projection: []string{"name", "id", "address.line1"}}
	result := command.Run()
	minifyResult, _ := util.MinifyJsonObject(result)

	if result == nil || string(minifyMock) != string(minifyResult) {
		t.Errorf("record must exists")
	}
}

func TestSelectWithoutProjection(t *testing.T) {
	mockData, _ := util.ToJsonObject(string(`[{
		"a5bdf0df-432e-4272-69ea-87b13ba06a8a":
		{
			"address": {
				"line1": "shankargarh",
				"state": "uttar-pradesh"
			},
			"dob": "15-05-1986",
			"email": "pooja.1234@gmail.com",
			"lastName": "chaturvedi",
			"name": "pooja"
		}
	}]`))
	minifyMock, _ := util.MinifyJsonObject(mockData)
	os.Setenv("DB_DATA", "./")
	filter := service.Filter{Table: "student_test", FilterKey: "name", FilterValue: "pooja"}
	command := service.Select{Filter: filter, Projection: []string{}}
	result := command.Run()
	minifyResult, _ := util.MinifyJsonObject(result)

	if result == nil || string(minifyMock) != string(minifyResult) {
		t.Errorf("record must exists")
	}
}

func TestSelectFromEmptyJsonTable(t *testing.T) {
	tableName := "empty_file"
	os.Setenv("DB_DATA", "./")
	jsonDb := util.JsonDb{Table: tableName}
	jsonDb.Truncate()
	filter := service.Filter{Table: tableName, FilterKey: "name", FilterValue: "pooja"}
	command := service.Select{Filter: filter, Projection: []string{}}
	result := command.Run()
	minifyResult, _ := util.MinifyJsonObject(result)

	if string(minifyResult) != "null" {
		t.Errorf("record must Not exists")
	}
}
