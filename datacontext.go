package orm

import (
	"github.com/eaciit/database/base"
	//"github.com/eaciit/errorlib"
)

type DataContext struct {
	//Adapter base.IAdapter
	Connection base.IConnection
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
