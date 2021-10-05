APP_NAME ?= flowerss
CONSUMER_NAME ?= consumer
IMAGE_NAME ?= flowerss


test:
	go test ./... -v

build: get
	go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go
	go build -o ./bin/$(CONSUMER_NAME) ./cmd/$(CONSUMER_NAME)/main.go

get:
	go mod download

run:
	go run .

clean:
	rm flowerss-bot

rund:
	docker-compose up -d

logd:
	docker-compose logs -f

downd:
	docker-compose down --remove-orphans

image:
	docker build -t $(IMAGE_NAME) .
