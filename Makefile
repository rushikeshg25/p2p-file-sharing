PHONY: build

build:
	go build -o ./bin/p2p-share

clean:
	rm ./bin/p2p-share