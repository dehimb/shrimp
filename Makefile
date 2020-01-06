build:
	# compile proto files
	protoc --proto_path=proto --go_out=plugins=grpc:proto proto/port/port.proto
	# install dependencies
	go mod download
	# build portdomainservice
	GOOS=linux GOARCH=amd64 go build -o portdomainservice/portdomainservice portdomainservice/cmd/main.go
	# build clientapi
	GOOS=linux GOARCH=amd64 go build -o clientapi/clientapi clientapi/cmd/main.go
test:
	go test `go list ./... | grep -v -e "proto" -e "cmd"`
up: 
	docker-compose up -d --build
down:
	docker-compose stop
