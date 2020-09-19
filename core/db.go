package core

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zfd81/rooster/rsql"
	"github.com/zfd81/sunny/meta"
)

type ParrotDB struct {
	namespace string
	Name      string
	*rsql.DB
}

func (d *ParrotDB) GetNamespace() string {
	if d.namespace == "" {
		return meta.DefaultNamespace
	}
	return d.namespace
}

func NewDB(ds *meta.DataSource) (*ParrotDB, error) {
	var driverName, dsn string
	if strings.ToLower(ds.Driver) == "mysql" {
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", ds.User, ds.Password, ds.Host, ds.Port, ds.Database)
	}
	db, err := rsql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &ParrotDB{
		namespace: ds.Namespace,
		Name:      ds.Name,
		DB:        db,
	}, nil
}
