package office

import (
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
)

var _db *orm.DataContext

func SetDb(conn dbox.IConnection) error {
	CloseDb()
	_db = orm.New(conn)
	return nil
}

func CloseDb() {
	if _db != nil {
		_db.Close()
	}
}

func DB() *orm.DataContext {
	return _db
}
