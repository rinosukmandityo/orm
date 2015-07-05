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

func (d *DataContext) Register(m IModel) IModel {
	m.SetM(m)
	m.SetCtx(d)
	return m
}

func (d *DataContext) Close() {
	d.Connection.Close()
}

func (d *DataContext) SetConnectionFromConfig(name string) error {
	d.ConnectionName = name
	if d.ConnectionName == "" {
		d.ConnectionName = "Default"
	}

	connType := strings.ToLower(config.Get("Connection_" + d.ConnectionName + "_Default").(string))
	host := config.Get("Connection_" + d.ConnectionName + "_Host").(string)
	username := config.Get("Connection_" + d.ConnectionName + "_Username").(string)
	password := config.Get("Connection_" + d.ConnectionName + "_Password").(string)
	database := config.get("Connection_" + d.ConnectionName + "_Password").(string)

	if connType == "mongodb" {
		conn := mongodb.NewConnection(host, username, password, database)
		return nil
	} else {
		return errorlib.Error(packageName, modCtx, "SetConnectionFromConfig", "Connection for MySQL is not yet implemented")
	}
}
