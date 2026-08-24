run:
	go run ./cmd/main.go

build: linux windows mac

linux:
	mkdir -p dist/linux
	GOOS=linux GOARCH=amd64 go build -o dist/linux/pfin ./cmd/main.go

windows:
	mkdir -p dist/win
	GOOS=windows GOARCH=amd64 go build -o dist/win/pfin.exe ./cmd/main.go

mac:
	mkdir -p dist/osx
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/pfin ./cmd/main.go

clean-dist:
	rm -rf ./dist 2>/dev/null
