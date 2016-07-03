package ormpackage

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/ormgen"
)

type ORMObject struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            string `bson:"_id" json:"_id"`
	Enable        bool
}

func (o *ORMObject) TableName() string {
	return "ormobjects"
}

func NewORMObject() *ORMObject {
	o := new(ORMObject)
	return o
}

func ORMObjectFind(filter *dbox.Filter, fields string, limit, skip int) dbox.ICursor {
	config := makeFindConfig(fields, skip, limit)
	if filter != nil {
		config.Set("where", filter)
	}
	c, _ := DB().Find(new(ORMObject), config)
	return c
}

func ORMObjectGet(id interface{}) *ORMObject {
	emp := new(ORMObject)
	e := DB().GetById(emp, id)
	if e != nil {
		return nil
	}
	return emp
}

func ORMObjectFindByEnable(enable bool,
	fields string, limit, skip int) dbox.ICursor {
	c := ORMObjectFind(dbox.Eq("enable", enable),
		fields, limit, skip)
	return c
}
