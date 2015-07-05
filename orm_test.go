package orm

import (
	//"github.com/eaciit/database/base"
	"fmt"
	"github.com/eaciit/database/mongodb"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

type UserModel struct {
	ModelBase `bson:"-"`
	Id        string `bson:"_id"`
	FullName  string `bson:"fullname"`
	Email     string
	Password  string
	Enable    int `bson:"enable"`
}

var e error

func (u *UserModel) Init() *UserModel {
	u.M = u
	return u
}

func prepareContext() (*DataContext, error) {
	conn := mongodb.NewConnection("localhost:27017", "", "", "ectest")
	if eConnect := conn.Connect(); eConnect != nil {
		return nil, eConnect
	}
	ctx := NewDataContext(conn)
	return ctx, nil
}

func (u *UserModel) TableName() string {
	return "ORMUsers"
}

func TestInsert(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	count := 1000000
	wg := sync.WaitGroup{}
	wg.Add(count)
	ctx, e := prepareContext()
	defer ctx.Close()
	t0 := time.Now()
	for a := 0; a < 100; a++ {
		go func(a int, wg *sync.WaitGroup) {
			var u *UserModel
			for x := 0; x < 10000; x++ {
				i := a*10000 + x + 1
				wg.Done()
				u = new(UserModel)
				ctx.Register(u)
				u.Id = "user" + strconv.Itoa(i)
				u.FullName = "ORM User " + strconv.Itoa(i)
				u.Email = "ormuser01@email.com"
				u.Password = "mbahmu kepet"
				u.Enable = 1
				e = u.Save()
				if e != nil {
					t.Errorf("Error Load: %s", e.Error())
					//return
				}
				fmt.Println("Inserted: " + strconv.Itoa(i))
			}
		}(a, &wg)
	}

	defer func() {
		d0 := time.Since(t0)
		fmt.Printf("Completed in %v \n", d0)
	}()

	wg.Wait()
}

/*
func TestInsert(t *testing.T) {
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()

	u := new(UserModel)
	ctx.Register(u)
	u.Id = "user01"
	u.FullName = "ORM User 01"
	u.Email = "ormuser01@email.com"
	u.Password = "mbahmu kepet"
	u.Enable = 1
	e = u.Save()
	if e != nil {
		t.Errorf("Error Load: %s", e.Error())
		return
	}
}

func TestLoad(t *testing.T) {
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()

	u := new(UserModel)
	ctx.Register(u)
	e = u.GetById("user01")
	if e != nil {
		t.Errorf("Error Load: %s", e.Error())
		return
	} else {
		fmt.Printf("UserModel: %v \n", u)
		fmt.Println("")
	}
}

func TestDelete(t *testing.T) {
	ctx, e := prepareContext()
	if e != nil {
		t.Errorf("Error Connect: %s", e.Error())
		return
	}
	defer ctx.Close()
	u := new(UserModel)
	ctx.Register(u)
	e = u.GetById("user01")
	if e == nil {
		fmt.Printf("Will Delete UserModel: %v \n", u.M)
		e = u.Delete()
		if e != nil {
			t.Errorf("Error Load: %s", e.Error())
			return
		} else {
			fmt.Printf("UserModel: %v has been deleted \n", u)
			fmt.Println("")
		}
	}
}
*/
