.PHONY: extension


TREE_SITTER_SOURCE_DIR=./tree-sitter-refal5

extension:
	@code --extensionDevelopmentPath=$(pwd)/vscode
 
run:
	@go run ./cmd/server/main.go

build:
	@mkdir out
	@GOOS=linux GOARCH=arm64  go build -o ./out/server ./cmd/server/main.go 

generate-parser:
	@cd $(TREE_SITTER_SOURCE_DIR) && tree-sitter generate 
	@mkdir -p internal/tree_sitter_refal5/tree_sitter
	@cp -r $(TREE_SITTER_SOURCE_DIR)/src/tree_sitter/*  internal/tree_sitter_refal5/tree_sitter
	@cp $(TREE_SITTER_SOURCE_DIR)/src/parser.c internal/tree_sitter_refal5
	@cp $(TREE_SITTER_SOURCE_DIR)/src/scanner.c internal/tree_sitter_refal5

clean:
	@rm -rf ./out
