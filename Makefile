NAME=dextre

default: drain

build:
	go build -o output/$(NAME) .

drain: build
	./output/$(NAME) drain --node ip-172-20-100-230.eu-west-1.compute.internal --skip-validation=true

roll-node: build
	./output/$(NAME) roll nodes --role node
