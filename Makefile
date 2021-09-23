include common.mk

## build: Build app
.PHONY: build
build: build-backend

## build-backend: Build backend
.PHONY: build-backend
build-backend:
	$Q cd backend && make build

## lint: Run all lints
.PHONY: lint
lint: lint-commits lint-proto | ; $(info $(M) running linters…)

## lint-commits: Run commits linter
.PHONY: lint-commits
lint-commits: $(commitlint) | ; $(info $(M) running commit linters…)
	$Q ./scripts/commitlint.sh

## lint-proto: Run proto linter
.PHONY: lint-proto
lint-proto: $(buf) | ; $(info $(M) running proto linters…)
	$Q $(buf) lint

## bump-version: Bump version
.PHONY: bump-version
bump-version: $(standard-version) | ; $(info $(M) bumping app version…)
	$Q ./scripts/bump_version.sh

## proto-backend: Generate backed proto
PROTO_BACKEND_OUT = ./backend/gen/proto
.PHONY: proto-backend
proto-backend: $(protoc) $(protoc-gen-go) $(protoc-gen-validate) $(protoc-gen-twirp) $(gowrap) | ; $(info $(M) generate protos…)
	$Q rm -rf $(PROTO_BACKEND_OUT)
	$Q mkdir -p $(PROTO_BACKEND_OUT)
	$Q for file in $$(find ./proto -name '*.proto' -not -path "*validate*"); do \
		$(protoc) -I ./proto -I ./tools/protoc/include \
			--twirp_out=$(PROTO_BACKEND_OUT) \
			--twirp_opt=paths=source_relative \
			--go_out=$(PROTO_BACKEND_OUT) \
			--go_opt=paths=source_relative \
			--validate_out="lang=go:$(PROTO_BACKEND_OUT)" \
			--validate_opt=paths=source_relative \
			$$file; \
    done
	$Q cd ./backend && $(gowrap) gen -p ./gen/proto/filebrowser/v1 -i FileService -t twirp_validate -g -o ./gen/proto/filebrowser/v1/file_service.validate.go

## proto-ts: Generate ts rpc client proto files
PROTO_FRONTEND_OUT = ./frontend/src/gen/proto
.PHONY: proto-frontend
proto-frontend: $(protoc) $(protoc-gen-twirp_ts) $(protoc-gen-ts_proto)
	$Q rm -rf $(PROTO_FRONTEND_OUT)
	$Q mkdir -p $(PROTO_FRONTEND_OUT)
	$Q for file in $$(find ./proto -name '*.proto' -not -path "*validate*"); do \
		$(protoc) -I ./proto -I ./tools/protoc/include \
			--ts_proto_out=$(PROTO_FRONTEND_OUT) \
			--ts_proto_opt=esModuleInterop=true \
			--ts_proto_opt=env=browser \
			--ts_proto_opt=outputClientImpl=false \
			--twirp_ts_out=$(PROTO_FRONTEND_OUT) \
			--twirp_ts_opt="ts_proto" \
			--twirp_ts_opt=client_only \
			$$file; \
  	done

## clean: Remove generated files
.PHONY: clean
clean:
#	tools
	$Q rm -rf $(TOOLS_DIR)/bin
	$Q rm -rf $(TOOLS_DIR)/node_modules
#	frontend
	$Q rm -rf $(BASEPATH)/frontend/node_modules
	$Q rm -rf $(BASEPATH)/frontend/dist
