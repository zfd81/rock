package core

import (
	"fmt"
	"strings"

	"github.com/zfd81/rooster/types/container"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rooster/rsql"
)

type RockDB struct {
	namespace string
	Name      string
	db        *rsql.DB
}

func (d *RockDB) GetNamespace() string {
	if d.namespace == "" {
		return meta.DefaultNamespace
	}
	return d.namespace
}

func (d *RockDB) GetName() string {
	return d.Name
}

func (d *RockDB) QueryMap(query string, arg interface{}) (container.Map, error) {
	return d.db.QueryMap(query, arg)
}
func (d *RockDB) QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]container.Map, error) {
	return d.db.QueryMapList(query, arg, pageNumber, pageSize)
}

func (d *RockDB) Query(query string, arg interface{}) (*rsql.Rows, error) {
	return d.db.Query(query, arg)
}

func (d *RockDB) Exec(query string, arg interface{}) (int64, error) {
	return d.db.Exec(query, arg)
}

func (d *RockDB) Save(arg interface{}, table ...string) (int64, error) {
	return d.db.Save(arg, table...)
}

func (d *RockDB) BatchSave(arg []interface{}, table ...string) (int64, error) {
	return d.db.BatchSave(arg, table...)
}

func NewDB(ds *meta.DataSource) (*RockDB, error) {
	var driverName, dsn string
	if strings.ToLower(ds.Driver) == "mysql" {
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", ds.User, ds.Password, ds.Host, ds.Port, ds.Database)
	}
	db, err := rsql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &RockDB{
		namespace: ds.Namespace,
		Name:      ds.Name,
		db:        db,
	}, nil
}
