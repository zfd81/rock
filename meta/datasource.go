package meta

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
