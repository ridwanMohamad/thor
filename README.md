# Thor
Thor is a service build with go language, that handle user authentication, user management, role management, module & menu management, lookup value management

### Function Of This Service
- User authentication
- Session Generation
- Lookup Value
- Role management
- Menu management
- Role & Menu Matrix
- 
### Requirements
- Golang 1.17

### Installation
#### Local
```shell
//clone project
git clone https://github.com/Titipaja/thor.git
//change directory to project
cd thor
//downloading go module dependency
go mod tidy
```
#### Docker
```shell
//1ST TIME
git clone -b main git@github.com:Titipaja/thor.git
cd thor
sudo su
docker build -t titipaja/thor-server:1.0.0 .
docker run -d -p 9080:9080 --net bridge --name thor-service titipaja/thor-server:1.0.0

//2ND TIME & Continously
cd thor
git pull
sudo su
docker build -t titipaja/thor-server:1.0.1 .
docker kill thor-service
docker rm thor-service
docker run -d -p 9080:9080 --net bridge --name thor-service titipaja/thor-server:1.0.1
```
### Migration Files
- Not implemented yet

### Environment Variables
 ```
APPS_NAME=thor
APPS_VERSION=1.0.27
APPS_HTTP_PORT=9080

DB_MASTER_HOST=asgard-titip-626e.aivencloud.com
DB_MASTER_USERNAME=app_thor
DB_MASTER_PASSWORD=AVNS_6gEsWAwGG7KJpsG
DB_MASTER_NAME=defaultdb
DB_MASTER_SHCEMA=thor
DB_MASTER_PORT=18882
DB_MASTER_CONNECTION_TIMEOUT= 20s
DB_MASTER_IDLE_CONNECTION=5
DB_MASTER_MAX_OPEN_CONNECTION=10
DB_MASTER_DEBUG_MODE=true

MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_SENDER=denysetiawan28@gmail.com
MAIL_SENDER_NAME=TitipAja
MAIL_PASSWORD=qbjkdvaqurbjzhlm
MAIL_SMTP_AUTH=true
MAIL_START_TLS=true
 ```
