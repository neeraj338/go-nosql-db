package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-nosql-db/constants"
	"io/ioutil"
	"os"
	"strings"
)

type JsonDb struct {
	Table string
}

func GetBasePath() string {
	return getEnv("DB_DATA", "/tmp/")
}

func (db *JsonDb) FilePath() string {
	return fmt.Sprintf("%v%v.json", GetBasePath(), db.Table)
}

func (db *JsonDb) Save(json []byte) error {
	fileSize := db.GetFileSize()
	var formatForAppendJson []byte = json
	if fileSize <= 0 {
		err := ioutil.WriteFile(db.FilePath(), json, 0644)
		if err != nil {
			return err
		}
		return nil
	} else {
		jsonString := string(formatForAppendJson)
		trimPrefixBrace := strings.TrimPrefix(jsonString, "{")
		formatForAppendJson = []byte(fmt.Sprintf(",%v", trimPrefixBrace))
	}

	f, err := os.OpenFile(db.FilePath(), os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(-1, 2); err != nil {
		return err
	}
	if _, err := f.WriteAt(formatForAppendJson, fileSize-1); err != nil {
		return err
	}
	return nil
}

func (db *JsonDb) DeleteRecord(begin int64, recordLength int64, position constants.FilePosition) error {

	f, err := os.OpenFile(db.FilePath(), os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer f.Close()
	seekPos := map[constants.FilePosition]int64{constants.BeginingOfFile: begin, constants.MidOfFile: begin - int64(1), constants.EndOfFile: begin - int64(1)}[position]
	if _, err := f.Seek(seekPos, 0); err != nil {
		return err
	}

	blankString := strings.Repeat(" ", int(recordLength)+1)
	if _, err := f.WriteAt([]byte(blankString), seekPos); err != nil {
		return err
	}
	return nil
}

func (db *JsonDb) Truncate() error {
	f, err := os.OpenFile(db.FilePath(), os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %q for truncation: %v", db.Table, err)
	}
	defer f.Close()
	return nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func (db *JsonDb) IsDbTableExists() (bool, error) {
	_, err := os.Stat(db.FilePath())
	if err == nil {
		return true, err
	} else {
		return false, err
	}
}
func (db *JsonDb) GetFileSize() int64 {

	info, err := os.Stat(db.FilePath())

	if err != nil {
		return -1
	}
	x := info.Size()
	return x
}

func JsonPrettyPrint(data interface{}) string {
	byteArr, err := json.MarshalIndent(data, "", "\t")
	if err != nil || string(byteArr) == "null" {
		return ""
	}
	return string(byteArr)
}

func ToJsonObject(data string) (interface{}, error) {
	var doc interface{}
	err := json.Unmarshal([]byte(data), &doc)
	return doc, err
}

func MinifyJsonString(jsonString string) ([]byte, error) {
	cp := bytes.NewBuffer([]byte{})
	err := json.Compact(cp, []byte(jsonString))
	return cp.Bytes(), err
}

func MinifyJsonObject(jsonObject interface{}) ([]byte, error) {
	cp := bytes.NewBuffer([]byte{})
	jsonByteArr, marshalErr := json.Marshal(jsonObject)
	if marshalErr != nil {
		return cp.Bytes(), marshalErr
	}
	err := json.Compact(cp, jsonByteArr)
	return cp.Bytes(), err
}
