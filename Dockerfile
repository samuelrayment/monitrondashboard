FROM google/golang

WORKDIR /gopath/src/github.com/samuelrayment/monitrondashboard
ADD . /gopath/src/github.com/samuelrayment/monitrondashboard
RUN go get github.com/samuelrayment/monitrondashboard/...
RUN go install github.com/samuelrayment/monitrondashboard/...

CMD []
ENTRYPOINT ["/gopath/bin/monidash"]
