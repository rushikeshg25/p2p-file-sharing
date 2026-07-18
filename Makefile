.PHONY: build clean executable

build:
	mkdir -p ./bin
	go build -o ./bin/p2p-share

clean:
	rm -f ./bin/p2p-share

executable:
	chmod +x ./bin/p2p-share
