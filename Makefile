PHONY: build clean executable

build:
	go build -o ./bin/p2p-share

clean:
	rm ./bin/p2p-share

executable:
	chmod +x ./bin/p2p-share