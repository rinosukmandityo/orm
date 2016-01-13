package orm

import (
	"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/dbox"
	err "github.com/eaciit/errorlib"
	tk "github.com/eaciit/toolkit"
	"strings"
)

type DataContext struct {
	//Adapter dbox.IAdapter
	ConnectionName string
	Connection     dbox.IConnection

	pooling bool
	//adapters       map[string]dbox.IAdapter
}

func (d *DataContext) SetPooling(p bool) *DataContext {
	d.pooling = p
	return d
}

func (d *DataContext) Pooling() bool {
	return d.pooling
}

func New(conn dbox.IConnection) *DataContext {
	ctx := new(DataContext)
	ctx.Connection = conn
	//ctx.adapters = map[string]dbox.IAdapter{}
	return ctx
}

func NewFromConfig(name string) (*DataContext, error) {
	ctx := new(DataContext)
	//ctx.adapters = map[string]dbox.IAdapter{}
	eSet := ctx.setConnectionFromConfigFile(name)
	if eSet != nil {
		return ctx, eSet
	}
	return ctx, nil
}

func (d *DataContext) Find(m IModel, parms tk.M) dbox.ICursor {
	////_ = "breakpoint"
	q := d.Connection.NewQuery().From(m.TableName())
	if qe := parms.Get("where", nil); qe != nil {
		q = q.Where(qe.([]*dbox.Filter)...)
	}
	if qe := parms.Get("order", nil); qe != nil {
		q = q.Order(qe.([]string)...)
	}
	if qe := parms.Get("skip", nil); qe != nil {
		q = q.Skip(qe.(int))
	}
	if qe := parms.Get("limit", nil); qe != nil {
		q = q.Take(qe.(int))
	}
	//fmt.Printf("Debug Q: %s\n", tk.JsonString(q))
	c, _ := q.Cursor(nil)
	return c
}

func (d *DataContext) GetById(m IModel, id interface{}) (bool, error) {
	var e error
	q := d.Connection.NewQuery().SetConfig("pooling", d.Pooling()).From(m.(IModel).TableName()).Where(dbox.Eq("_id", id))
	c, _ := q.Cursor(nil)
	e = c.Fetch(m, 1, false)
	if e != nil {
		return false, err.Error(packageName, modCtx, "GetById", e.Error())
	}
	return true, nil
}

func (d *DataContext) Insert(m IModel) error {
	q := d.Connection.NewQuery().SetConfig("pooling", d.Pooling()).From(m.TableName()).Insert()
	e := q.Exec(tk.M{"data": m})
	return e
}

func (d *DataContext) Save(m IModel) error {
	var e error
	if m.RecordId() == nil {
		m.PrepareId()
	}
	if e = m.PreSave(); e != nil {
		return err.Error(packageName, modCtx, m.TableName()+".PreSave", e.Error())
	}
	q := d.Connection.NewQuery().SetConfig("pooling", d.Pooling()).From(m.TableName()).Save()
	e = q.Exec(tk.M{"data": m})
	if e = m.PostSave(); e != nil {
		return err.Error(packageName, modCtx, m.TableName()+",PostSave", e.Error())
	}
	return e
}

func (d *DataContext) Delete(m IModel) error {
	q := d.Connection.NewQuery().SetConfig("pooling", d.Pooling()).From(m.TableName()).Delete()
	e := q.Exec(tk.M{"data": m})
	return e
}

func (d *DataContext) DeleteMany(m IModel, where *dbox.Filter) error {
	var e error
	q := d.Connection.NewQuery().SetConfig("pooling", d.Pooling()).From(m.TableName()).Delete()
	if where != nil {
		q.Where(where)
	}
	e = q.Exec(tk.M{"where": where})
	return e
}

func (d *DataContext) Close() {
	d.Connection.Close()
}

func (d *DataContext) setConnectionFromConfigFile(name string) error {
	d.ConnectionName = name
	if d.ConnectionName == "" {
		d.ConnectionName = fmt.Sprintf("Default")
	}

	connType := strings.ToLower(config.Get("Connection_" + d.ConnectionName + "_Type").(string))
	host := config.Get("Connection_" + d.ConnectionName + "_Host").(string)
	username := config.Get("Connection_" + d.ConnectionName + "_Username").(string)
	password := config.Get("Connection_" + d.ConnectionName + "_Password").(string)
	database := config.Get("Connection_" + d.ConnectionName + "_database").(string)

	ci := new(dbox.ConnectionInfo)
	ci.Host = host
	ci.UserName = username
	ci.Password = password
	ci.Database = database

	conn, eConnect := dbox.NewConnection(connType, ci)
	if eConnect != nil {
		return err.Error(packageName, modCtx, "SetConnectionFromConfigFile", eConnect.Error())
	}
	if eConnect = conn.Connect(); eConnect != nil {
		return err.Error(packageName, modCtx, "SetConnectionFromConfigFile", eConnect.Error())
	}
	d.Connection = conn
	return nil
}
