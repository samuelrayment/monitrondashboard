#!/bin/sh

go get ./...
go install ./...
cp $GOPATH/bin/monidash monidash.tmp
sudo docker build -t monitrondashboard .
rm monidash.tmp
