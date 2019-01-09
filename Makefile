NAME=dextre

default: drain

build:
	go build -o output/$(NAME) .

drain: build
	./output/$(NAME) drain --node ip-172-20-120-231.eu-west-1.compute.internal

restart: build
	./output/$(NAME) restart --label app=feed --namespace dev
