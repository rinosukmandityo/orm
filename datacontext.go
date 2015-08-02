package orm

import (
	"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	err "github.com/eaciit/errorlib"
	"strings"
)

type DataContext struct {
	//Adapter base.IAdapter
	ConnectionName string
	Connection     base.IConnection
	adapters       map[string]base.IAdapter
}

func New(conn base.IConnection) *DataContext {
	ctx := new(DataContext)
	ctx.Connection = conn
	ctx.adapters = map[string]base.IAdapter{}
	return ctx
}

func NewFromConfig(name string) (*DataContext, error) {
	ctx := new(DataContext)
	ctx.adapters = map[string]base.IAdapter{}
	eSet := ctx.setConnectionFromConfigFile(name)
	if eSet != nil {
		return ctx, eSet
	}
	return ctx, nil
}

func (d *DataContext) Register(m IModel) IModel {
	if m.Ctx() != nil {
		return m
	}

	m.SetCtx(d)
	a, ok := d.adapters[m.TableName()]
	if !ok {
		a = d.Connection.Adapter(m.TableName())
		d.adapters[m.TableName()] = a
	}
	return m
}

func (d *DataContext) Find(m IModel, parms T) base.ICursor {
	_ = "breakpoint"
	return d.Connection.Table(m.TableName(), parms)
}

func (d *DataContext) GetById(m IModel, id interface{}) error {
	var e error
	//return err.Error(packageName, modModel, "GetById", err.NotYetImplemented)
	adapter, e := d.adapter(m)
	cursor, _, e := adapter.Run(base.DB_SELECT, nil, O{"find": O{"_id": id}})
	if e != nil {
		return err.Error(packageName, modCtx, "GetById", e.Error())
	}
	b, e := cursor.Fetch(m)
	if b == false {
		return err.Error(packageName, modCtx, "GetById", fmt.Sprintf("Record with id:%v could not be found", id))
	} else if e != nil {
		return err.Error(packageName, modCtx, "GetById", fmt.Sprintf("Error parse record with id:%v | %s", id, e.Error()))
	} else {
		m.SetCtx(d)
	}
	return nil
}

func (d *DataContext) Insert(m IModel) error {
	return d.saveOrInsert(m, base.DB_INSERT)
}

func (d *DataContext) Save(m IModel) error {
	return d.saveOrInsert(m, base.DB_SAVE)
}

func (d *DataContext) Delete(m IModel) error {
	a, e := d.adapter(m)
	if e != nil {
		return e
	}
	_, _, e = a.Run(base.DB_DELETE, m, nil)
	return e
}

func (d *DataContext) Close() {
	d.Connection.Close()
}

func (d *DataContext) saveOrInsert(m IModel, dbOp base.DB_OP) error {
	var e error
	a, e := d.adapter(m)
	if e != nil {
		return e
	}
	m.PrepareId()
	e = m.PreSave()
	if e != nil {
		return e
	}
	_, _, e = a.Run(dbOp, m, nil)
	e = m.PostSave()
	if e != nil {
		return e
	}
	return e
}

func (d *DataContext) adapter(m IModel) (base.IAdapter, error) {
	m = d.Register(m)
	tableName := m.TableName()
	_ = "breakpoint"
	a, ok := d.adapters[tableName]
	if !ok {
		return nil, err.Error(packageName, modCtx, "adapter", "Adapter "+tableName+" is not yet initialized")
	}
	return a, nil
}

func (d *DataContext) setConnectionFromConfigFile(name string) error {
	d.ConnectionName = name
	if d.ConnectionName == "" {
		d.ConnectionName = "Default"
	}

	connType := strings.ToLower(config.Get("Connection_" + d.ConnectionName + "_Type").(string))
	host := config.Get("Connection_" + d.ConnectionName + "_Host").(string)
	username := config.Get("Connection_" + d.ConnectionName + "_Username").(string)
	password := config.Get("Connection_" + d.ConnectionName + "_Password").(string)
	database := config.Get("Connection_" + d.ConnectionName + "_Database").(string)

	if connType == "mongodb" {
		conn := mongodb.NewConnection(host, username, password, database)
		if eConnect := conn.Connect(); eConnect == nil {
			d.Connection = conn
		} else {
			return err.Error(packageName, modCtx, "SetConnectionFromConfigFile", eConnect.Error())
		}
	} else {
		return err.Error(packageName, modCtx, "SetConnectionFromConfig", "Connection for "+connType+" is not yet implemented")
	}
	return nil
}
