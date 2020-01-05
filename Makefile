build:
	# compile proto files
	protoc --proto_path=proto --go_out=plugins=grpc:proto proto/port/port.proto
	# install dependencies
	go mod download
	# build port-domain-service
	GOOS=linux GOARCH=amd64 go build -o port-domain-service/port-domain-service port-domain-service/cmd/main.go
	# build client-api
	GOOS=linux GOARCH=amd64 go build -o client-api/client-api client-api/cmd/main.go

up: 
	docker-compose up -d --build

down:
	docker-compose stop
