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
	ModelBase  `bson:"-",json:"-"`
	ID         string `bson:"_id",json:"_id"`
	FullName   string `bson:"fullname"`
	Age        int
	Email      string
	Password   string
	RandomDate time.Time
	Enable     int `bson:"enable"`
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

func (u *UserModel) RecordID() interface{} {
	return u.ID
}

var ctx *DataContext

func TestInsert(t *testing.T) {
	//t.Skip()
	ctx, _ := prepareContext()
	defer ctx.Close()

	ctx.DeleteMany(new(UserModel), nil)

	t0 := time.Now()
	count := 20
	for i := 1; i <= count; i++ {
		fmt.Printf("Insert user no %d ...", i)
		u := new(UserModel)
		u.ID = "user" + strconv.Itoa(i)
		u.FullName = "ORM User " + strconv.Itoa(i)
		u.Age = tk.RandInt(20) + 20
		u.Email = "ormuser01@email.com"
		u.Password = "mbahmu kepet"
		u.Enable = 1
		u.RandomDate = time.Now().Add(time.Duration(int64(tk.RandInt(1000)) * int64(time.Minute)))
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
		u := new(UserModel)
		u.ID = fmt.Sprintf("user%d", i)
		fmt.Printf("Update user %s ...", u.ID)
		e := ctx.GetById(u, u.ID)
		//e := ctx.GetById(u, "user3")
		if e == nil {
			u.FullName = "ORM User X" + strconv.Itoa(i)
			u.Email = "ormuser01@email.com"
			u.Password = "mbahmu kepet tha ?"
			u.Enable = 0
			e = ctx.Save(u)
			if e != nil {
				t.Errorf("Error Load %d: %s", i, e.Error())
				return
			} else {
				fmt.Println("OK")
			}
		} else {
			fmt.Println("NOK ..." + e.Error())
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
	e = ctx.GetById(u, "user2")
	if e == nil {
		fmt.Printf("Will Delete UserModel:\n %s \n", tk.JsonString(u))
		e = ctx.Delete(u)
		if e != nil {
			t.Errorf("Error Load: %s", e.Error())
			return
		} else {
			tk.Unjson(tk.Jsonify(u), u)
			fmt.Printf("UserModel: %v has been deleted \n", u.RandomDate.UTC())
			fmt.Println("")
		}
	} else {
		t.Errorf("Delete error: %s", e.Error())
	}
}

func TestLoadAll(t *testing.T) {
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()

	tk.Println("Test Load All")
	c, _ := ctx.Find(new(UserModel), tk.M{
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
		users := []UserModel{}
		if count > 0 {
			e = c.Fetch(&users, 0, false)
		}
		if e == nil {
			fmt.Printf("OK...")
			fmt.Printf("Record(s) found: %d\nSample of first record: %s \n", count, tk.IfEq(count, 0, "", users[0].Email))
			fmt.Println("")
		} else {
			fmt.Println("NOK")
			t.Error(e.Error())
		}
	}
}
