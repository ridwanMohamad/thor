FROM golang:1.19-alpine as builder

#get git
RUN apk update && \
apk upgrade && \
apk add --no-cache git && \
apk add --no-cache tzdata

#add user for golang and maintainer
#RUN addgroup -S golang && adduser -S golang -G golang
#USER golang:golang
MAINTAINER titipaja.id

#working directory
ADD . /opt/apps
WORKDIR /opt/apps

#copy resource
COPY . .

#building
#RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o thor-server

#FROM scratch
FROM alpine:3.16

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /opt/apps/resources /opt/apps/resources

COPY --from=builder /opt/apps/thor-server /opt/apps/thor-server

WORKDIR /opt/apps

#expose network
EXPOSE 9080

#running
CMD [ "/opt/apps/thor-server" ]