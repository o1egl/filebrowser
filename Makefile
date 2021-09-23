include common.mk

## build: Build app
.PHONY: build
build: build-backend

## build-backend: Build backend
.PHONY: build-backend
build-backend:
	$Q cd backend && make build

## proto: Compile proto files
.PHONY: proto
proto: $(protoc) $(protoc-gen-go) $(protoc-gen-validate) $(protoc-gen-twirp) $(gowrap) | ; $(info $(M) generate protos…)
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

## proto-ts: Generate ts rpc client proto files
TS_OUT = "./frontend/my-app/gen/proto"
.PHONY: proto-ts
proto-ts: $(protoc) $(protoc-gen-twirp_ts) $(protoc-gen-ts_proto)
	$Q rm -rf $(TS_OUT)
	$Q mkdir -p $(TS_OUT)
	$Q for file in $$(find ./proto -name '*.proto' -not -path "*github.com*"); do \
  		$(protoc) -I ./proto -I ./tools/protoc/include \
  			--ts_proto_out=$(TS_OUT) \
  			--ts_proto_opt=esModuleInterop=true \
  			--ts_proto_opt=env=browser \
            --ts_proto_opt=outputClientImpl=false \
  			--twirp_ts_out=$(TS_OUT) \
  			--twirp_ts_opt="ts_proto" \
  			--twirp_ts_opt=client_only \
  			$$file; \
  	done

## lint: Run all lints
.PHONY: lint
lint: lint-commits | ; $(info $(M) running linters…)

## lint-commits: Run commit linters
.PHONY: lint-commits
lint-commits: $(commitlint) | ; $(info $(M) running commit linters…)
	$Q ./scripts/commitlint.sh

## bump-version: Bump version
.PHONY: bump-version
bump-version: $(standard-version) | ; $(info $(M) bumping app version…)
	$Q ./scripts/bump_version.sh

## clean: Remove generated files
.PHONY: clean
clean:
#	tools
	$Q rm -rf $(TOOLS_DIR)/bin
	$Q rm -rf $(TOOLS_DIR)/node_modules
#	frontend
	$Q rm -rf $(BASEPATH)/frontend/node_modules
	$Q rm -rf $(BASEPATH)/frontend/dist
