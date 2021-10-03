build: build-linux build-mac build-windows

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -ldflags "-w -s" -o ./bin/resolver-linux main.go
build-mac:
	CGO_ENABLED=0 GOOS=darwin go build -mod=readonly -a -ldflags "-w -s" -o ./bin/resolver-mac main.go
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/resolver.exe main.go