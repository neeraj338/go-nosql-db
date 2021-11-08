package service

import (
	"go-nosql-db/util"

	uuid "github.com/nu7hatch/gouuid"
)

type Persist struct {
	Table string
	Data  string
}

func (p *Persist) Run() interface{} {

	jsonObjet, err := util.ToJsonObject(p.Data)
	if err != nil {
		panic(err)
	}
	uid, _ := uuid.NewV4()
	var assignId map[string]interface{} = make(map[string]interface{})
	assignId[uid.String()] = jsonObjet

	jsonToSave, _ := util.MinifyJsonObject(assignId)

	db := &util.JsonDb{Table: p.Table}
	db.Save(jsonToSave)
	return assignId
}
