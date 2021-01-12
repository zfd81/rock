module github.com/zfd81/rock

go 1.14

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/dop251/goja v0.0.0-20201221183957-6b6d5e2b5d80
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/golang/protobuf v1.4.1
	github.com/google/uuid v1.1.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/zfd81/rooster v0.0.0-20200509130230-5f2b9d86cd8a
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	google.golang.org/grpc v1.27.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.4.3
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
)
