package main

var baseGo string = `
package {0}

import (
	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/ormgen"
	"github.com/eaciit/toolkit"
)

var _dbs map[string]*orm.DataContext

func initDbs() {
	if _dbs == nil {
		_dbs = map[string]*orm.DataContext{}
	}
}

func SetDb(conn dbox.IConnection, ids ...string) error {
	initDbs()
	CloseDb(ids...)
	dbID := "default"
	if len(ids) > 0 {
		dbID = ids[0]
	}
	var (
		_db *orm.DataContext
		e   bool
	)
	if _db, e = _dbs[dbID]; !e {
		_db = orm.New(conn)
	}
	_dbs[dbID] = _db
	return nil
}

func CloseDb(ids ...string) {
	initDbs()
	dbID := "default"
	if len(ids) > 0 {
		dbID = ids[0]
	}
	if _db, e := _dbs[dbID]; e {
		if _db != nil {
			_db.Close()
		}
	}
}

func DB(ids ...string) *orm.DataContext {
	initDbs()
	dbID := "default"
	if len(ids) > 0 {
		dbID = ids[0]
	}
	_db := _dbs[dbID]
	return _db
}

func makeFindConfig(fields string, skip, limit int) toolkit.M {
	config := toolkit.M{}
	if fields != "" {
		fieldses := strings.Split(fields, ",")
		config.Set("select", fieldses)
	}
	if skip > 0 {
		config.Set("skip", skip)
	}
	if limit > 0 {
		config.Set("take", limit)
	}
	return config
}
`
