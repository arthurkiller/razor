#!/bin/bash

go get github.com/arthurkiller/razor/demo/client
go get github.com/arthurkiller/razor/demo/server

server &
sleep 1
client

echo "All test have done"
