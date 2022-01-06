FROM golang:1.17

WORKDIR /go/src/app

COPY . .

RUN make dist-x86_64

RUN mv ./extensions /opt

RUN ls -la /opt/*