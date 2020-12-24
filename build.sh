#! /bin/bash

#macOS版本
go build -o bin/macos/rock ./main.go
go build -o bin/macos/rockctl ./rockctl/main.go

#Linux版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/rock ./main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/rockctl ./rockctl/main.go

#Windows版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/rock.exe ./main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/rockctl.exe ./rockctl/main.go