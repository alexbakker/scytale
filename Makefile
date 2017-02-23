all: scycli scyserver

scycli: prep
	go build -o build/bin/scycli github.com/Impyy/scytale/cmd/scycli

scyserver: prep scyserver-assets
	go build -o build/bin/scyserver github.com/Impyy/scytale/cmd/scyserver

scyserver-assets:
	go run vendor/github.com/Impyy/go-embed/*.go -pkg=main -input=cmd/scyserver/assets -output=cmd/scyserver/assets.go

prep:
	mkdir -p build/bin

clean:
	rm -rf build
	rm -f cmd/scyserver/assets.go
