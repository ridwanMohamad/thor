FROM alpine:3.16
#get git
RUN apk update && \
apk upgrade && \
apk add --no-cache git && \
apk add --no-cache tzdata

MAINTAINER titipaja.id

#working directory
ADD ./resources /opt/apps/resources
ADD ./thor-server /opt/apps/thor-server

WORKDIR /opt/apps

EXPOSE 9080

#running
CMD [ "/opt/apps/thor-server" ]