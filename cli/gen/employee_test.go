package office

import (
	"log"
	"strconv"
	"testing"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm/cli/gen"
)

func InitCall() {
	conn, _ := dbox.NewConnection("mongo", &dbox.ConnectionInfo{"localhost:27017", "ormdb", "", "", nil})
	err := conn.Connect()
	if err != nil {
		log.Printf("CONN ERR %+v \n", err)
	}
	log.Printf("CONN %+v \n", conn)
	SetDb(conn)
}
func TestSave(t *testing.T) {
	InitCall()
	e := NewEmployee()
	e.ID = "emp" + strconv.Itoa(10)
	e.Title = "TEST Title"
	e.Address = "SOme Address"
	log.Printf("DB %+v", DB())
	log.Printf("e %+v", e)
	DB().Save(e)
}
func TestFindById(t *testing.T) {
	InitCall()
	e := EmployeeGetByID("emp10")
	log.Printf("EMP => %+v\n", e)
}
func TestFindByEnable(t *testing.T) {
	InitCall()
	emps := EmployeeFindByEnable(true, []string{"_id"}, 0, 0)
	log.Printf("EMPS => %+v\n", emps)
}
func TestFindByTitle(t *testing.T) {
	InitCall()
	emps := EmployeeFindByTitle("TEST Title", []string{"_id"}, 0, 0)
	log.Printf("EMPS => %+v\n", emps)
}
