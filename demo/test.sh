#!/bin/bash

cd ./client && go install
cd ./server && go install

server &
client
