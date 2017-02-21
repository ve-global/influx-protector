# Stolen from the nice folks at hashicorp: https://github.com/hashicorp/terraform/blob/master/Makefile

TEST?=./rules
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go')

default: test vet

deps:
	godep restore

# bin generates the releaseable binaries
bin: deps fmt test
	@TF_RELEASE=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: fmt
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

quickdev:
	@TF_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# test runs the unit tests
test: fmt
	go test $(TEST) $(TESTARGS)

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: bin default test vet fmt
