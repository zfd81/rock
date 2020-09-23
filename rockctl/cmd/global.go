package cmd

import "github.com/zfd81/rock/http"

type GlobalFlags struct {
	Endpoints []string
	User      string
	Password  string
}

var (
	client = http.New()
)

func url(path string) string {
	return "http://" +
		globalFlags.Endpoints[0] +
		"/rock/" +
		path
}
