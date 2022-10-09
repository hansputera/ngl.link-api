FROM golang:alpine3.16

RUN apk update

RUN mkdir -p /home/nglapi
COPY . /home/nglapi
WORKDIR /home/nglapi

RUN go mod download
RUN go build

CMD ["./nglapi"]
