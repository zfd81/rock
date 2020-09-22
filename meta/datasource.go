package meta

import "encoding/json"

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

func (ds *DataSource) Key() string {
	return DataSourceKey(ds.Namespace, ds.Name)
}

func (ds *DataSource) String() (string, error) {
	bytes, err := json.Marshal(ds)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewDataSource(jsonstr []byte) (*DataSource, error) {
	ds := &DataSource{}
	if err := json.Unmarshal(jsonstr, ds); err != nil {
		return nil, err
	}
	return ds, nil
}
