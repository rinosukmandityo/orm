package orm

import (
	"github.com/eaciit/config"
	"github.com/eaciit/database/base"
	"github.com/eaciit/database/mongodb"
	"github.com/eaciit/errorlib"
	"strings"
)

type DataContext struct {
	//Adapter base.IAdapter
	ConnectionName string
	Connection     base.IConnection
}

func NewDataContext(conn base.IConnection) *DataContext {
	ctx := new(DataContext)
	ctx.Connection = conn
	return ctx
}

func NewDataContextFromConfig(name string) (*DataContext, error) {
	ctx := new(DataContext)
	eSet := ctx.setConnectionFromConfigFile(name)
	if eSet != nil {
		return ctx, eSet
	}
	return ctx, nil
}

func (d *DataContext) Register(m IModel) IModel {
	m.SetM(m)
	m.SetCtx(d)
	return m
}

func (d *DataContext) Insert(m IModel) error {
	m = d.Register(m)
	return m.Insert()
}

func (d *DataContext) Save(m IModel) error {
	m = d.Register(m)
	return m.Save()
}

func (d *DataContext) Delete(m IModel) error {
	m = d.Register(m)
	return m.Delete()
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
			return errorlib.Error(packageName, modCtx, "SetConnectionFromConfigFile", eConnect.Error())
		}
	} else {
		return errorlib.Error(packageName, modCtx, "SetConnectionFromConfig", "Connection for "+connType+" is not yet implemented")
	}
	return nil
}
