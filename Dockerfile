FROM golang:1.16-alpine AS build

WORKDIR /build

COPY go.* ./

RUN go mod download

COPY tjcounter.go .

RUN CGO_ENABLED=0 go build -o tjcounter .

WORKDIR /

RUN cp /build/tjcounter .

COPY templates /templates/

EXPOSE 8181

ENTRYPOINT ["/tjcounter"]