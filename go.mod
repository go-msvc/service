module github.com/go-msvc/service

go 1.12

require (
	github.com/apex/log v1.1.4
	github.com/go-msvc/config v0.0.0-20191007220741-e8517c202db8
	github.com/go-msvc/log v0.0.0-20200515104948-e039d1c2f30d
	github.com/pkg/errors v0.8.1
)

replace github.com/go-msvc/config => ../config

replace github.com/go-msvc/errors => ../errors

replace github.com/go-msvc/log => ../log
