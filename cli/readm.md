#ormgen
ormgen is a cli command to automatically generate GO files for model based on a text file

##usage
```
ormgen -file=pathtofile -out=outputfolder
```

pathtofile will define path of file to be read. If no parameter is being defined, by default it will find for default.orm on working directory

ouputfolder will define location of generated GO file. If no parameter is being defined, by default it will be working directory

File created will have name convention into xxx.go where xxx is struct name converted into lower case

##sample of default.orm
```
package office 

/* This is a remark */
/* C:Employee Commented remark - any commented remak started with C: will be copied over to generated code and eliminiate its C: part*/
struct Employee     /*Create employee.go*/
TableName:employeTables /*Tablename on orm will be employeTables ... if no tablename define default is employees (plural name of struct in lower case) */
ID:string:primary_key:bson`_id:json`_id /*with primary_key exist, this will generate getByID function by default*/
Title:string:default_EMPTY TITLE:json`title /*Example of default value of string */ 
Enable:bool:default_true    /*Field enable, type is bool, default value when New is true*/
GetByID()          /*Will generate GetByID(id string) *Employee */
FindByTitle()       /*Will generate FindByTitle(title string) dbox.ICursor */
FindByEnable()   /*Will genderate FindByEnable(enable bool) dbox.ICursor */
/*Every 'FindBy' with field name structure, orm will find the type of parameter*/

/*C:Department this is a commented remark and should be copied over to code for documentation purpose*/
struct Department
ID:string:primary_key
Title:string
Enable:bool:defaut_true
OwnerID:string:reference_Employee /*Field EmployeeID, is a reference to Employee. ormgen should automatically created func (d *Department) Owner()*Employee */

For example please view sample.orm file
```

Generated file should overwrite existing file and should reserve any changes that has been made outside any definition create within .orm file 
### Limitation ##
currently only supporting for function and variable changes made by user, any other changes will be overwritten
####


Above .orm file should generate base.go, employee.go and department.go and 

Sample of base.go
```go
package office
import(
    "github.com/eaciit/dbox"
	"github.com/eaciit/orm"
)

var _db *orm.DataContext

func SetDb(conn dbox.IConnection)error{
    CloseDb()
    _db = orm.New(conn)
    return nil
}

func CloseDb(){
    if _db!=nil{
        _db.Close()    
    }
}

func DB() *orm.DataContext{
    return _db
}
```

Sample or employee.go will be
```go
package office

import (
	"github.com/eaciit/dbox"
	. "github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Employee struct {
	ModelBase `bson:"-",json:"-"`
	ID        string ` bson:"_id" , json:"_id" `
	Title     string `json:"title" `
	Address   string ` bson:"address" `
	Enable    bool
	LastLogin time.Time
	OtherId   bson.ObjectId
}

func EmployeeGetByID(id string) *Employee {
	employee := new(Employee)
	DB().GetById(employee, id)
	return employee
}
func EmployeeFindByTitle(title string, order []string, skip, limit int) dbox.ICursor {
	c, _ := DB().Find(new(Employee),
		toolkit.M{}.Set("where", dbox.Eq("title", title)).
			Set("order", order).
			Set("skip", skip).
			Set("limit", limit))
	return dbox.NewCursor(c)
}

func EmployeeFindByEnable(enable bool, order []string, skip, limit int) dbox.ICursor {
	c, _ := DB().Find(new(Employee),
		toolkit.M{}.Set("where", dbox.Eq("enable", enable)).
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
```

Sample or department.go will be
```go

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

func (e *Department) Owner() *Employee {
	return EmployeeGetByID(e.OwnerID)
}
func NewDepartment() *Department {
	e := new(Department)
	e.Enable = true
	return e
}

func (e *Department) TableName() string {
	return "departments"
}
```