compile: compile-arm64 compile-amd64
	lipo -create -output evn-pilot-conversion ./evn-pilot-conversion-arm64 ./evn-pilot-conversion-amd64

compile-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./evn-pilot-conversion-arm64 ./cmd/

compile-amd64:
	GOOS=darwin GOARCH=amd64 go build -o ./evn-pilot-conversion-amd64 ./cmd/
