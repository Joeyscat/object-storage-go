#!/usr/bin/env bash

#kill $(ps -ef|grep go-build|awk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')

kill $(ps -ef | grep go-build | awk '$0 !~/grep/ {print $2}' | tr -s '\n' ' ')
kill $(ps -ef | grep server.go | awk '$0 !~/grep/ {print $2}' | tr -s '\n' ' ')
