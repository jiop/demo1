# DEMO 1

## Purpose of this project

Micro-service usage in web application with auth and webserver separated in two different service using docker.

## Usage

`make` generate docker images for all different services

`docker-compose up` start all containers for this demo with networks, volumes, etc.

## Web Server

### Description

Responsible of handling requests from users.

Serves html file in order to be interpreted by a web browser.

Answers with appropriate http code.

### Routes

POST /login (with username and password)

POST /logout (with cookie containing username and token)

GET /protected-content (with cookie containing username and token)

GET /unprotected-content (no parameters)

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

POST /login (with username and password)

POST /logout (with username and token)

### Networking

Port 8002 is open for the web server.

Share docker network `demo1net` with the Web Server.

### System

Uses Alpine or Golang docker image. Probably golang image for development and alpine with compiled golang project for production.

Share one volume to store the codebase during development (folder /app on docker container)

## TODO

### General

- Switch to a golang image during development
- Create real images for authserver and webserver for deploy process
- Add the possibility to scale auth server and web server.

### Web Server

- Choose a web framework in golang
- main.go launch the webserver and handle requests
- conf.yml store IP address of the AuthServer
- make use of cookies

### Auth Server

- choose a different framework in golang
- main.go launch the authserver and handle requests
- main.go store (hardcoded for now) the list of all the couples username and password
- main.go store all couples username and token

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
