package office

import (
	. "github.com/eaciit/orm"
)

type Department struct {
	ModelBase `bson:"-",json:"-"`
	ID        string ` bson:"_id" , json:"_id" `
	Title     string
	Enable    bool
	OwnerID   string
}

func (e *Department) RecordID() interface{} {
	return e.ID
}
func NewDepartment() *Department {
	e := new(Department)
	e.Enable = true
	return e
}
func (e *Department) TableName() string {
	return "departments"
}
