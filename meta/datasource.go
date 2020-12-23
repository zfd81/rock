package meta

import (
	"gopkg.in/mgo.v2/bson"
)

type DataSource struct {
	Namespace string `json:"namespace"` //命名空间 注:不能包含"/"
	Name      string `json:"name"`
	Driver    string `json:"driver"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
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
