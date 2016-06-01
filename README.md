# DEMO 1

## Avant-propos

Mostly inspired by an article about test integration by semaphoreci.

TODO: add url

Due to my limited access to internet (limited 3g access), I do not use any web framework or special packages in this project.

## Purpose of this project

Micro-service usage in web application with auth and webserver separated in two different service using docker.

## Usage

generate docker images for all different services
```
make
```

build and start all containers for this demo with networks, volumes, etc.
```
docker-compose up --build
```

Define an environment variable containing the docker machine ip.
```
DOCKER_MACHINE_IP=$(docker-machine default ip)
```

Login with username "user2" and password "pass2"
```
curl -XPOST $DOCKER_MACHINE_IP:8080/login -d username=user2 -d password=pass2 -c cookie.txt
```

Use stored cookies to access protected content
```
curl -XPOST $DOCKER_MACHINE_IP:8080/protected-content -c cookie.txt -b cookie.txt
```

## Web Server

### Description

Responsible of handling requests from users.

Serves html file in order to be interpreted by a web browser.

Answers with appropriate http code.

Serves a cookie containing username and password for further access to protected content.

### Routes

POST /login (with username and password, responds with json and cookies)

POST /logout (with cookie containing username and token, responds json and empty cookies)

GET /protected-content (with cookie containing username and token, responds json)

GET /unprotected-content (no parameters, responds with json)

### Networking

Share docker network `demo1net` with the Auth Server.

### System

Uses Alpine or Golang docker image. Probably golang image for development and alpine with compiled golang project for production.

Share one volume to store the codebase during development (folder /app on docker container)

## Auth Server

### Description

Responsible of login and logout users.

Serves json data in order to be interpreted by the web server.

Answers with appropriate http code.

Stores usernames, passwords and tokens.

### Routes

POST /login (with username and password, responds with json)

POST /logout (with username and token, responds with json)

### Networking

Port 8002 is open for the web server.

Share docker network `demo1net` with the Web Server.

### System

Uses Alpine or Golang docker image. Probably golang image for development and alpine with compiled golang project for production.

Share one volume to store the codebase during development (folder /app on docker container)

## TODO

### General

- Switch to a golang image during development => Dockerfile.dev now use golang image
- Create real images for authserver and webserver for deploy process => need to create a Dockerfile.prod and Dockerfile.test
- Add the possibility to scale auth server and web server.

### Web Server

- Choose a web framework in golang
- main.go launch the webserver and handle requests => DONE
- conf.yml store IP address of the AuthServer
- make use of cookies => DONE

### Auth Server

- choose a different framework in golang
- main.go launch the authserver and handle requests
- main.go store (hardcoded for now) the list of all the couples username and password => DONE
- use Redis or another database to store username and password

## Notes

### Build the Golang compiler docker image

cd ../gobuilder && docker build -t jiop/gobuilder .

### Build the Golang executables

cd ../authserver && docker run --rm -it -v $PWD/src:/src/ -v $PWD/build:/build/ jiop/gobuilder

cd ../webserver && docker run --rm -it -v $PWD/src:/src/ -v $PWD/build:/build/ jiop/gobuilder

### Launch each server

cd ../authserver && docker build -t jiop/authserver .

cd ../authserver && docker run -it --rm -v $PWD/build:/app jiop/authserver

cd ../webserver && docker build -t jiop/webserver .

cd ../webserver && docker run -it --rm -v $PWD/build:/app jiop/webserver

### Other stuff (text)

Receive request with username/password or username/token. In the first case, it responds with a token and a 200 http code. In the second case, it answer 200 http code when the user is authenticated or a forbiden http code if not.

Receive request from the user (authenticated or not). Responds to /login and /logout requests.  Moreover, it responds to /protected-content and /unprotected-content. The first one is accessible only after a successful request to /login and the second one, does not need any login.

docker network create demo1net

docker run -it --rm --name webserver1 -P --net demo1net jiop/webserver

docker run -it --rm --name authserver1 --net demo1net jiop/authserver
