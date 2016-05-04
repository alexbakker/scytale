all: scytale-cli scytale-server

scytale-cli: prep
	go build -o build/bin/scytale-cli github.com/Impyy/Scytale/cmd/scytale-cli

scytale-server: prep scytale-server-assets
	go build -o build/bin/scytale-server github.com/Impyy/Scytale/cmd/scytale-server

scytale-server-assets:
	go run vendor/github.com/Impyy/go-embed/*.go -pkg=main -input=cmd/scytale-server/assets -output=cmd/scytale-server/assets.go

prep:
	mkdir -p build/bin

clean:
	rm -rf build
	rm -f cmd/scytale-server/assets.go
