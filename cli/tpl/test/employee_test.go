package test

import (
	"log"
	"math"
	"strconv"
	"testing"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/ormgen/cli/tpl"
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
	log.Printf("CONN %+v \n", conn)
	ormpackage.SetDb(conn)
}

func TestSave(t *testing.T) {
	InitCall()
	for i := 1; i <= 1000; i++ {
		e := ormpackage.NewORMObject()
		e.ID = "emp" + strconv.Itoa(i)
		//e.Title = toolkit.Sprintf("Test Title %d", i)
		//e.Address = toolkit.Sprintf("Address %d", i)
		//e.LastLogin = time.Now()
		if math.Mod(float64(i), 2) == 0 {
			e.Enable = true
		} else {
			e.Enable = false
		}
		//log.Printf("DB %+v", ormpackage.DB())
		//log.Printf("e %+v", toolkit.JsonString(e))
		ormpackage.DB().Save(e)
	}
}

func TestFindById(t *testing.T) {
	//InitCall()
	e := ormpackage.ORMObjectGet("emp10")
	log.Printf("EMP => %+v\n", toolkit.JsonString(e))
}

func TestFindByEnable(t *testing.T) {
	//InitCall()
	emps := ormpackage.ORMObjectFindByEnable(true, "_id,title,enable", 0, 0)
	defer emps.Close()
	log.Printf("EMPS => %+v\n", emps.Count())

	i := 0
	for {
		emp := ormpackage.NewORMObject()
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
	ormpackage.DB().Close()
}
