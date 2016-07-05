package office

import (
	"time"

	"github.com/eaciit/dbox"
)

type Employee struct {
	ID      string `bson:"_id"`
	Title   string
	Created time.Time
	Enable  bool
}

func (o *Employee) TableName() string {
	return "employees"
}
func NewEmployee() *Employee {
	o := new(Employee)
	o.Enable = true
	return o
}
func EmployeeFind(filter *dbox.Filter, fields string, limit, skip int) dbox.ICursor {
	config := makeFindConfig(fields, skip, limit)
	if filter != nil {
		config.Set("where", filter)
	}
	c, _ := DB().Find(new(Employee), config)
	return c
}
