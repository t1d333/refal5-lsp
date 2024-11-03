.PHONY: extension

extension:
	@code --extensionDevelopmentPath=$(pwd)/vscode
 
run:
	@go run ./cmd/server/main.go

build:
	@mkdir out
	@GOOS=linux GOARCH=arm64  go build -o ./out/server ./cmd/server/main.go 

clean:
	@rm -rf ./out
