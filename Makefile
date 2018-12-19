PROJ_PATH = $(GOPATH)/src/github.com/aaronflower/shippy-service-consignment
build:
	protoc -I. --go_out=plugins=micro:$(PROJ_PATH) proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build -o service.consignment
	docker build --rm -t service.consignment .

run:
	docker run -p 50051:50051  \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns service.consignment
	
clean:
	go clean
	rm service.consignment
