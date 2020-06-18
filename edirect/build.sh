#!/bin/sh

go mod init edirect
go mod tidy
go build xtract.go common.go
go build rchive.go common.go
go build j2x.go
go build t2x.go
