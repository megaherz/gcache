#!/usr/bin/env bash
go build
./gcache -addr=:8080  & ./gcache -psw=123 -addr=:8081 && fg
rm gcache