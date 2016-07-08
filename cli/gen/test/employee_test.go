package test

import (
	"log"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/ormgen/cli/gen"
	"github.com/eaciit/toolkit"
)

func InitCall() {
	conn, _ := dbox.NewConnection("mongo",
		&dbox.ConnectionInfo{"localhost:27123",
			"ormdb", "", "", nil})
	err := conn.Connect()
	if err != nil {
		log.Printf("CONN ERR %+v \n", err)
	}
	//log.Printf("CONN %+v \n", conn)
	office.SetDb(conn)
}

func TestSave(t *testing.T) {
	InitCall()
	for i := 1; i <= 1000; i++ {
		e := office.NewEmployee()
		e.ID = "emp" + strconv.Itoa(i)
		e.Title = toolkit.Sprintf("Test Title %d", i)
		e.Address = toolkit.Sprintf("Address %d", i)
		e.LastLogin = time.Now()
		if math.Mod(float64(i), 2) == 0 {
			e.Enable = true
		} else {
			e.Enable = false
		}
		//log.Printf("DB %+v", office.DB())
		//log.Printf("e %+v", toolkit.JsonString(e))
		office.DB().Save(e)
	}
}

func TestGetById(t *testing.T) {
	//InitCall()
	emp, e := office.EmployeeGetByID("emp110", "")
	if e != nil {
		t.Fatal(e.Error())
	}
	log.Printf("EMP => %+v\n", toolkit.JsonString(emp))
}

func TestGetByTitleEnable(t *testing.T) {
	//InitCall()
	emp, e := office.EmployeeGetByTitleEnable("Test Title 116", true, "")
	if e != nil {
		t.Fatal(e.Error())
	}
	log.Printf("EMP => %+v\n", toolkit.JsonString(emp))
}

func TestFindByEnable(t *testing.T) {
	//InitCall()

	emp6, _ := office.EmployeeGetByID("emp6", "")
	emp8, _ := office.EmployeeGetByID("emp8", "")
	office.DB().Delete(emp6)
	office.DB().Delete(emp8)

	emps := office.EmployeeFindByEnable(true,
		"_id,title,enable", 0, 0)
	defer emps.Close()
	log.Printf("EMPS => %+v\n", emps.Count())

	i := 0
	for {
		emp := office.NewEmployee()
		e := emps.Fetch(emp, 1, false)
		if e != nil {
			break
		}
		toolkit.Println(toolkit.JsonString(emp))
		i++
		if i == 10 {
			break
		}
	}
}

func TestClose(t *testing.T) {
	office.DB().Close()
}
