#!/bin/bash

go get github.com/arthurkiller/razor/demo/client
go get github.com/arthurkiller/razor/demo/server

server &
client
