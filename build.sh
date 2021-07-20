#!/bin/sh
GOOS=windows GOARCH=amd64 go build -o mgit.exe
go build -o mgit
