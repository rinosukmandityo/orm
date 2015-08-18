package orm

import (
	//"fmt"
	"github.com/eaciit/config"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	err "github.com/eaciit/errorlib"
	tk "github.com/eaciit/toolkit"
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

func (d *DataContext) Find(m IModel, parms tk.M) base.ICursor {
	//_ = "breakpoint"
	q := d.Connection.Query().From(m.TableName())
	if qe := parms.Get("where", nil); qe != nil {
		q = q.Where(qe.(*base.QE))
	}
	if qe := parms.Get("order", nil); qe != nil {
		q = q.OrderBy(qe.([]string)...)
	}
	if qe := parms.Get("skip", nil); qe != nil {
		q = q.Skip(qe.(int))
	}
	if qe := parms.Get("limit", nil); qe != nil {
		q = q.Limit(qe.(int))
	}
	return q.Cursor(nil)
}

func (d *DataContext) GetById(m IModel, id interface{}) error {
	return nil
}

func (d *DataContext) Insert(m IModel) error {
	q := d.Connection.Query().From(m.TableName()).Insert()
	_, _, e := q.Run(tk.M{"data": m})
	return e
	//return d.saveOrInsert(m, base.DB_INSERT)
}

func (d *DataContext) Save(m IModel) error {
	q := d.Connection.Query().From(m.TableName()).Save()
	_, _, e := q.Run(tk.M{"data": m})
	return e
}

func (d *DataContext) Delete(m IModel) error {
	q := d.Connection.Query().From(m.TableName()).Delete()
	_, _, e := q.Run(tk.M{"data": m})
	return e
}

func (d *DataContext) DeleteMany(m IModel, where *base.QE) error {
	var e error
	q := d.Connection.Query().From(m.TableName()).Delete()
	if where == nil {
		_, _, e = q.Run(nil)
	} else {
		_, _, e = q.Run(tk.M{"where": where})
	}
	return e
}

func (d *DataContext) Close() {
	d.Connection.Close()
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
