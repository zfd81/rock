package meta

import (
	"gopkg.in/mgo.v2/bson"
)

type DataSource struct {
	Namespace string `yaml:"namespace"` //命名空间 注:不能包含"/"
	Name      string `yaml:"name"`
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
}

func (ds *DataSource) EtcdKey() string {
	return DataSourceEtcdKey(ds.Namespace, ds.Name)
}

func (ds *DataSource) String() (string, error) {
	bytes, err := bson.Marshal(ds)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewDataSource(bytes []byte) (*DataSource, error) {
	ds := &DataSource{}
	if err := bson.Unmarshal(bytes, ds); err != nil {
		return nil, err
	}
	return ds, nil
}
