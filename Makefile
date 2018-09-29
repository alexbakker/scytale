export GO15VENDOREXPERIMENT=1

all: scycli scyserver

scycli: prep
	go build -o build/bin/scycli github.com/alexbakker/scytale/cmd/scycli

scyserver: prep
	go build -o build/bin/scyserver github.com/alexbakker/scytale/cmd/scyserver

prep:
	mkdir -p build/bin

clean:
	rm -rf build
