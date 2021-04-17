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
