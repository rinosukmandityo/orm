package office

import (
	"github.com/eaciit/dbox"
	. "github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

type Employee struct {
	ModelBase `bson:"-",json:"-"`
	ID        string ` bson:"_id" , json:"_id" `
	Title     string `json:"title" `
	Address   string ` bson:"address" `
	Enable    bool
}

func EmployeeGetByID(id string) *Employee {
	employee := new(Employee)
	DB().GetById(employee, id)
	return employee
}
func EmployeeFindByTitle(title string, order []string, skip, limit int) dbox.ICursor {
	c, _ := DB().Find(new(Employee),
		toolkit.M{}.Set("where", []*dbox.Filter{dbox.Eq("title", title)}).
			Set("order", order).
			Set("skip", skip).
			Set("limit", limit))
	return dbox.NewCursor(c)
}
func EmployeeFindByEnable(enable bool, order []string, skip, limit int) dbox.ICursor {
	c, _ := DB().Find(new(Employee),
		toolkit.M{}.Set("where", []*dbox.Filter{dbox.Eq("enable", enable)}).
			Set("order", order).
			Set("skip", skip).
			Set("limit", limit))
	return dbox.NewCursor(c)
}
func (e *Employee) RecordID() interface{} {
	return e.ID
}
func NewEmployee() *Employee {
	e := new(Employee)
	e.Title = "EMPTY TITLE"
	e.Enable = true
	return e
}
func (e *Employee) TableName() string {
	return "employeTables"
}
