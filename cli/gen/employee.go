package office

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	"time"
)

type Employee struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            string `bson:"_id"`
	Title         string
	Address       string
	LastLogin     time.Time
	Created       time.Time
	Enable        bool
}

func (o *Employee) TableName() string {
	return "employees"
}
func NewEmployee() *Employee {
	o := new(Employee)
	o.Enable = true
	return o
}
func EmployeeFind(filter *dbox.Filter, fields, orders string, limit, skip int) dbox.ICursor {
	config := makeFindConfig(fields, orders, skip, limit)
	if filter != nil {
		config.Set("where", filter)
	}
	c, _ := DB().Find(new(Employee), config)
	return c
}
func EmployeeGet(filter *dbox.Filter, orders string, skip int) (emp *Employee, err error) {
	config := makeFindConfig("", orders, skip, 1)
	if filter != nil {
		config.Set("where", filter)
	}
	c, ecursor := DB().Find(new(Employee), config)
	if ecursor != nil {
		return nil, ecursor
	}
	defer c.Close()

	emp = new(Employee)
	err = c.Fetch(emp, 1, false)
	return emp, err
}

func EmployeeGetByID(pID string, orders string) (*Employee, error) {
	return EmployeeGet(dbox.Eq("_id", pID), "", 0)
}

func EmployeeGetByTitleEnable(pTitle string, pEnable bool, orders string) (*Employee, error) {
	return EmployeeGet(dbox.And(dbox.Eq("title", pTitle), dbox.Eq("enable", pEnable)), "", 0)
}

func EmployeeFindByEnable(pEnable bool, fields string, limit, skip int) dbox.ICursor {
	return EmployeeFind(dbox.Eq("enable", pEnable), "", "", 0, 0)
}
