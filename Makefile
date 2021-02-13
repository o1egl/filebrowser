include common.mk

## lint: Run all lints
.PHONY: lint
lint: lint-proto lint-commits | ; $(info $(M) running linters…)

## lint-proto: Run proto linters
.PHONY: lint-proto
lint-proto: install-buf | ; $(info $(M) running proto linters…)
	$Q $(buf) check lint

## lint-commits: Run commit linters
.PHONY: lint-commits
lint-commits: install-commitlint | ; $(info $(M) running commit linters…)
	$Q ./scripts/commitlint.sh

## release: Bump applivation version
.PHONY: bump-version
bump-version: install-standard-version | ; $(info $(M) bumping app version…)
	$Q ./scripts/bump_version.sh

## proto: Compile proto files
.PHONY: proto
proto: install-protoc install-protoc-gen-validate install-protoc-gen-twirp install-gowrap | ; $(info $(M) generate protos…)
	$Q rm -rf ./backend/gen/proto
	$Q mkdir -p ./backend/gen/proto
	$Q for file in $$(find ./proto -name '*.proto' -not -path "*github.com*"); do \
		$(protoc) -I ./proto -I ./tools/protoc/include \
			--twirp_out=./backend/gen/proto --twirp_opt=paths=source_relative \
			--go_out=./backend/gen/proto --go_opt=paths=source_relative \
			--validate_out="lang=go:./backend/gen/proto" --validate_opt=paths=source_relative \
			 $$file; \
    done
	$Q cd ./backend && $(gowrap) gen -p ./gen/proto/file/v1 -i FileService -t twirp_validate -g -o ./gen/proto/file/v1/file_service.validate.go
	$Q cd ./backend && $(gowrap) gen -p ./gen/proto/user/v1 -i UserService -t twirp_validate -g -o ./gen/proto/user/v1/user_service.validate.go

## clean: Remove generated files
.PHONY: clean
clean:
#	tools
	$Q rm -rf $(TOOLS_DIR)/bin
	$Q rm -rf $(TOOLS_DIR)/node_modules
	$Q rm -rf $(TOOLS_DIR)/protoc
#	frontend
	$Q rm -rf $(BASEPATH)/frontend/node_modules
	$Q rm -rf $(BASEPATH)/frontend/dist
