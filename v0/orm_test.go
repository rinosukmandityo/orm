package orm

import (
	"fmt"
	_ "github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	_ "github.com/eaciit/parallel"
	tk "github.com/eaciit/toolkit"
	"runtime"
	"strconv"
	_ "sync"
	"testing"
	"time"
)

type UserModel struct {
	ModelBase `bson:"-",json:"-"`
	Id        string `bson:"_id",json:"_id"`
	FullName  string `bson:"fullname"`
	Email     string
	Password  string
	Enable    int `bson:"enable"`
}

var e error

func (u *UserModel) Init() *UserModel {
	//u.M = u
	return u
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func prepareContext() (*DataContext, error) {
	conn := mongodb.NewConnection("localhost:27123", "", "", "ectest")
	if eConnect := conn.Connect(); eConnect != nil {
		return nil, eConnect
	}
	ctx := New(conn)
	return ctx, nil
}

func (u *UserModel) TableName() string {
	return "ORMUsers"
}

var ctx *DataContext

func TestLoadAll(t *testing.T) {
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()

	fmt.Println("Test Load All")
	c := ctx.Find(new(UserModel), tk.M{
		"where": nil,
		"order": []string{"_id"},
		"take":  0,
		"limit": 0,
	})
	defer c.Close()
	if c == nil {
		t.Errorf("Error Load: %s", c.Error)
		return
	} else {
		count := c.Count()
		user := new(UserModel)
		if count > 0 {
			_, e = c.Fetch(&user)
		}
		if e == nil {
			fmt.Printf("OK...")
			fmt.Printf("Record(s) found: %d\nSample of first record%v \n", count, tk.IfEq(count, 0, nil, user))
			fmt.Println("")
		} else {
			fmt.Println("NOK")
			t.Error(e.Error())
		}
	}
}

func TestInsert(t *testing.T) {
	//t.Skip()
	ctx, _ := prepareContext()
	defer ctx.Close()

	ctx.DeleteMany(new(UserModel), nil)

	t0 := time.Now()
	count := 10000
	for i := 0; i < count; i++ {
		fmt.Printf("Insert user no %d ...", i)
		u := new(UserModel)
		u.Id = "user" + strconv.Itoa(i)
		u.FullName = "ORM User " + strconv.Itoa(i)
		u.Email = "ormuser01@email.com"
		u.Password = "mbahmu kepet"
		u.Enable = 1
		e = ctx.Insert(u)
		if e != nil {
			t.Errorf("Error Load %d: %s", i, e.Error())
			return
		} else {
			fmt.Println("OK")
		}
	}
	fmt.Printf("Run process for %v \n", time.Since(t0))
}

func TestUpdate(t *testing.T) {
	t.Skip()
	ctx, _ := prepareContext()
	defer ctx.Close()

	t0 := time.Now()
	count := 100
	for i := 0; i < count; i++ {
		fmt.Printf("Insert user no %d ...", i)
		u := new(UserModel)
		u.Id = "user" + strconv.Itoa(i)
		u.FullName = "ORM User X" + strconv.Itoa(i)
		u.Email = "ormuser01@email.com"
		u.Password = "mbahmu kepet tha ?"
		u.Enable = 1
		e = ctx.Save(u)
		if e != nil {
			t.Errorf("Error Load %d: %s", i, e.Error())
			return
		} else {
			fmt.Println("OK")
		}
	}
	fmt.Printf("Run process for %v \n", time.Since(t0))
}

func TestDelete(t *testing.T) {
	t.Skip()
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()
	u := new(UserModel)
	e = ctx.GetById(u, "user01")
	if e == nil {
		fmt.Printf("Will Delete UserModel: %v \n", u)
		e = ctx.Delete(u)
		if e != nil {
			t.Errorf("Error Load: %s", e.Error())
			return
		} else {
			fmt.Printf("UserModel: %v has been deleted \n", u)
			fmt.Println("")
		}
	}
}
