FROM debian:buster AS ssl-generator
RUN apt-get update && \
    apt-get install -y openssl

RUN mkdir /certs
WORKDIR /certs
RUN openssl req \
    -new \
    -newkey rsa:4096 \
    -days 365 \
    -nodes \
    -x509 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com" \
    -keyout /certs/godock.key \
    -out /certs/godock.crt

FROM golang:1.13-buster AS build-env

RUN mkdir /builder
WORKDIR /builder
RUN git clone https://github.com/asciifaceman/godock.git
WORKDIR /builder/godock
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/godock
RUN chmod +x /go/bin/godock

FROM scratch
LABEL maintainer="Charles Corbett <nafredy@gmail.com>"
COPY --from=build-env /go/bin/godock /godock
COPY --from=ssl-generator /certs/godock.crt /godock.crt
COPY --from=ssl-generator /certs/godock.key /godock.key
ENV GODOCKCRT=/godock.crt
ENV GODOCKKEY=/godock.key
ENTRYPOINT [ "/godock" ]