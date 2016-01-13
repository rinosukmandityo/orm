package ormtest

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm/v1"
	tk "github.com/eaciit/toolkit"
	"strconv"
	"testing"
	"time"
)

type UserModel struct {
	ModelBase `bson:"-",json:"-"`
	Id        string `bson:"_id",json:"_id"`
	FullName  string `bson:"fullname"`
	Age       int
	Email     string
	Password  string
	Enable    int `bson:"enable"`
}

var e error

func (u *UserModel) Init() *UserModel {
	//u.M = u
	return u
}

func prepareContext() (*DataContext, error) {
	conn, _ := dbox.NewConnection("mongo", &dbox.ConnectionInfo{"localhost:27123", "ectest", "", "", nil})
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

	tk.Println("Test Load All")
	c := ctx.Find(new(UserModel), tk.M{
		"where": nil,
		"order": []string{"_id"},
		"take":  0,
		"limit": 0,
	})
	defer c.Close()

	if c == nil {
		t.Errorf("Error Load: Unable to init cursor")
		return
	} else {
		count := c.Count()
		user := new(UserModel)
		if count > 0 {
			_, e = c.Fetch(&user, 1, false)
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
	count := 100
	for i := 1; i <= count; i++ {
		fmt.Printf("Insert user no %d ...", i)
		u := new(UserModel)
		u.Id = "user" + strconv.Itoa(i)
		u.FullName = "ORM User " + strconv.Itoa(i)
		u.Age = tk.RandInt(20) + 20
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
	//t.Skip()
	ctx, _ := prepareContext()
	defer ctx.Close()

	t0 := time.Now()
	count := 10
	for i := 0; i < count; i++ {
		fmt.Printf("Update user no %d ...", i)
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
	//t.Skip()
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()
	u := new(UserModel)
	_, e = ctx.GetById(u, "user1")
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