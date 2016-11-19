#!/usr/bin/env bash
go build
./run -addr=:8080  & ./run -psw=123 -addr=:8081 && fg
rm run
