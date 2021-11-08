package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-nosql-db/constants"
	"go-nosql-db/util"
	"os"
	"regexp"
	"strings"

	"github.com/bcicen/jstream"
	"github.com/thedevsaddam/gojsonq"
)

type Filter struct {
	FilterKey   string
	FilterValue string
	Table       string
}

func (f *Filter) isFilterByValue() bool {
	if (len(strings.TrimSpace(f.FilterKey)) == 0) && (len(strings.TrimSpace(f.FilterValue)) > 0) {
		return true
	}
	return false
}
func (f *Filter) isFilterByKeyAndValue() bool {
	if (len(strings.TrimSpace(f.FilterKey)) > 0) && (len(strings.TrimSpace(f.FilterValue)) > 0) {
		return true
	}
	return false
}

type RecordPositions struct {
	AllRecordsSelectedFromTable bool
	Positions                   []*Position
}

func BuildRecordPosition(isAllrecordsSelectedFromTable bool, positios []*Position) RecordPositions {
	return RecordPositions{AllRecordsSelectedFromTable: isAllrecordsSelectedFromTable, Positions: positios}
}

type Position struct {
	Begin        int64
	Length       int64
	FilePosition constants.FilePosition
}

func BuildPosition(begining int64, length int64) *Position {
	return &Position{Begin: begining, Length: length, FilePosition: constants.EndOfFile}
}
func (rPos *Position) setPositionToBegining() {
	rPos.FilePosition = constants.BeginingOfFile
}
func (rPos *Position) setPositionToMiddle() {
	rPos.FilePosition = constants.MidOfFile
}

func (f *Filter) Search() ([]map[string]interface{}, RecordPositions) {
	dbTable := util.JsonDb{Table: f.Table}
	if _, err := dbTable.IsDbTableExists(); err != nil {
		panic(errors.New(fmt.Sprintf("DB Table (%v) not found", dbTable.FilePath())))
	}
	file, _ := os.Open(dbTable.FilePath())
	defer file.Close()
	var filterResult []map[string]interface{}
	var positions []*Position = make([]*Position, 0)
	decoder := jstream.NewDecoder(file, 1).EmitKV() // extract JSON values at a depth level of 1

	var firstIteration bool = true
	var isAllRecordsSelectedFromTable bool = true
	var currentIterationMv *jstream.MetaValue = nil
	var previousIterationMv *jstream.MetaValue = nil

	for mv := range decoder.Stream() {

		if len(positions) > 0 {
			if p := positions[0]; p.FilePosition == constants.BeginingOfFile && previousIterationMv == nil {
				p.Length = int64(mv.Offset) - p.Begin - 1
			}
			if p := positions[len(positions)-1]; p.FilePosition != constants.BeginingOfFile {
				p.setPositionToMiddle()
			}
		}

		previousIterationMv = currentIterationMv
		currentIterationMv = mv

		jsonByteArr, _ := json.Marshal(mv.Value)

		var p map[string]interface{}
		_ = json.Unmarshal(jsonByteArr, &p)
		uuid := p["key"].(string)

		var jMap map[string]interface{} = make(map[string]interface{})
		jMap[uuid] = p["value"]

		valBytes, _ := json.Marshal(jMap)
		jsonString := string(valBytes)

		isMatched := f.matchForFilterCondition(uuid, jsonString)

		if isMatched {
			filterResult = append(filterResult, jMap)
			recordPosition := BuildPosition(int64(mv.Offset), int64(mv.Length))
			if firstIteration {
				recordPosition.setPositionToBegining()
			} else {
				begin := recordPosition.Begin
				length := recordPosition.Length
				recordPosition.Begin = int64(previousIterationMv.Offset) + int64(previousIterationMv.Length) + 1
				recordPosition.Length = length + (begin - recordPosition.Begin)
			}

			positions = append(positions, recordPosition)
		} else {
			isAllRecordsSelectedFromTable = false
		}

		firstIteration = false
	}
	if firstIteration {
		isAllRecordsSelectedFromTable = false
	}

	return filterResult, BuildRecordPosition(isAllRecordsSelectedFromTable, positions)
}

func (f *Filter) matchForFilterCondition(uuid string, jsonString string) bool {
	var matchPresent = false
	if f.isFilterByValue() {
		r, _ := regexp.Compile(f.FilterValue)
		if r.MatchString(jsonString) || strings.EqualFold(f.FilterValue, uuid) {
			matchPresent = true
		}
	} else if f.isFilterByKeyAndValue() {
		jq := gojsonq.New().FromString(jsonString)
		res := jq.Find(fmt.Sprintf("%v.%v", uuid, f.FilterKey))

		if res != nil && res == f.FilterValue {
			matchPresent = true
		}
	}
	return matchPresent
}
