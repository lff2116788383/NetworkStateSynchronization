#!/bin/sh
ulimit -c unlimited
go build  -o ../bin/Server ./main/main.go
