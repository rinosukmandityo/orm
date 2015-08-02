package orm

import (
	//"github.com/eaciit/database/base"
	"fmt"
	"github.com/eaciit/database/mongodb"
	"github.com/eaciit/parallel"
	"runtime"
	"strconv"
	_ "sync"
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

func TestInsert(t *testing.T) {
	ctx, _ := prepareContext()
	defer ctx.Close()
	t0 := time.Now()
	count := 50000
	for i := 0; i < count; i++ {
		u := new(UserModel)
		ctx.Register(u)
		u.Id = "user" + strconv.Itoa(i)
		u.FullName = "ORM User " + strconv.Itoa(i)
		u.Email = "ormuser01@email.com"
		u.Password = "mbahmu kepet"
		u.Enable = 1
		e = ctx.Save(u)
		if e != nil {
			t.Errorf("Error Load: %s", e.Error())
			//return
		}
	}
	fmt.Printf("Run process for %v \n", time.Since(t0))
}

func TestInsertParallel(t *testing.T) {
	ctx, _ := prepareContext()
	defer ctx.Close()
	count := 50000
	workerCount := 100
	keys := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		keys = append(keys, i)
	}

	r := parallel.RunParallelJob(keys, workerCount, insertJob, parallel.T{"ctx": ctx})
	fmt.Printf("Run process for %v \n", r.Duration)
	if r.Status != "OK" {
		fmt.Printf("Error: %d fails \n %v \n", r.Fail, r.FailMessages)
	}
}

func insertJob(key interface{}, in parallel.T, result *parallel.JobResult) error {
	ctx := in["ctx"].(*DataContext)
	var u *UserModel
	i := key.(int)
	u = new(UserModel)
	//ctx.Register(u)
	u.Id = "user" + strconv.Itoa(i)
	u.FullName = "ORM User " + strconv.Itoa(i)
	u.Email = fmt.Sprintf("ormuser%d@email.com", i)
	u.Password = "mbahmu kepet " + strconv.Itoa(i)
	u.Enable = 1
	e = ctx.Save(u)
	if e != nil {
		//t.Errorf("Error Load: %s", e.Error())
		return fmt.Errorf("Error Load: %s", e.Error())
	}
	//fmt.Println("Insert : " + strconv.Itoa(i))
	return nil
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
