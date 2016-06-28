package office

import (
	"github.com/eaciit/orm"
    "github.com/eaciit/dbox"
    "github.com/eaciit/toolkit"
)

/*** OrmGen Auto Generate Code - Start ***/
type Employee struct {
	orm.ModelBase
	/**{struct_fields}**/
	ID    string `bson:"_id" json:"_id"`
	Title string
}

func (o *Employee) TableName() string{
    return "employees"
}

func EmployeeFind(filter *dbox.Filter)dbox.ICursor{
    config := toolkit.M{}
    c, _ := DB().Find(new(Employee), config)
    return c
}

/*** OrmGen Auto Generate Code - End ***/
