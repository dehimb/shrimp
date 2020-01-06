# Disclaimer
This 🦐  repo is for test purpose

## Requirements
To run this project, you need install git, go-1.13, docker, docker-compose, make, protoc, gcc

## How to use
Clone this repo, inside you find make file.<br/>
`make build` will build project<br/>
`make test` will run tests<br/>
`make up` will create containers and run them<br/>
`make down` will gracefull stop services<br/>
After start, you can send GET request to port :8080 to retrive information.<br/>
For example `host:8080/getPort/INHYD`
