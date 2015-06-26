package orm

import (
	"fmt"
	"github.com/eaciit/database/base"
	err "github.com/eaciit/errorlib"
)

type O map[string]interface{}

type IModel interface {
	GetById(interface{}) error
	PreSave() error
	PostSave() error
	SetM(IModel) IModel
	SetCtx(*DataContext)
	TableName() string
}

type ModelBase struct {
	M       IModel        `bson:"-"`
	ctx     *DataContext  `bson:"-"`
	adapter base.IAdapter `bson:"-"`
	//Id      interface{}   `bson:"_id"`
	//Title   string        `bson:omitempty`
}

func (m *ModelBase) SetM(md IModel) IModel {
	m.M = md
	return m
}

func (m *ModelBase) SetCtx(dc *DataContext) {
	m.ctx = dc
	tableName := m.TableName()
	m.adapter = dc.Connection.Adapter(tableName)
}

func (m *ModelBase) GetById(id interface{}) error {
	if m.ctx == nil {
		return err.Error(packageName, modCtx, "GetById", "Database Context is not yet initialized")
	}
	//return err.Error(packageName, modModel, "GetById", err.NotYetImplemented)
	cursor, _, e := m.adapter.Run(base.DB_SELECT, nil, O{"find": O{"_id": id}})
	if e != nil {
		return err.Error(packageName, modModel, "GetById", e.Error())
	}
	oldCtx := m.ctx
	oldAdapter := m.adapter
	oldM := m.M
	b, e := cursor.Fetch(m.M)
	fmt.Printf("Record: %v \n", m.M)
	if b == false {
		return err.Error(packageName, modCtx, "GetById", fmt.Sprintf("Record with id:%v could not be found", id))
	} else if e != nil {
		return err.Error(packageName, modCtx, "GetById", fmt.Sprintf("Error parse record with id:%v | %s", id, e.Error()))
	} else {
		m.ctx = oldCtx
		m.adapter = oldAdapter
		m.M = oldM
	}
	return nil
}

func (m *ModelBase) Insert() error {
	var e error
	e = m.PreSave()
	if e != nil {
		return e
	}
	_, _, e = m.adapter.Run(base.DB_INSERT, m.M, nil)
	if e != nil {
		return e
	}
	e = m.PostSave()
	if e != nil {
		return e
	}
	return nil
}

func (m *ModelBase) Save() error {
	var e error
	e = m.PreSave()
	if e != nil {
		return e
	}
	_, _, e = m.adapter.Run(base.DB_SAVE, m.M, nil)
	if e != nil {
		return e
	}
	e = m.PostSave()
	if e != nil {
		return e
	}
	return nil
}

func (m *ModelBase) Delete() error {
	var e error
	fmt.Printf("Value now: %v \n", m)
	_, _, e = m.adapter.Run(base.DB_DELETE, m.M, nil)
	if e != nil {
		return e
	}
	return nil
}

func (m *ModelBase) TableName() string {
	if m.M == nil {
		return "GenericTables"
	}
	return m.M.TableName()
}

func (m *ModelBase) PreSave() error {
	return nil
}

func (m *ModelBase) PostSave() error {
	return nil
}
