.DEFAULT_GOAL := default

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmtcheck errcheck vet
	@find ./cmd/* -maxdepth 1 -type d -exec go install "{}" \;

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

package: build
	@find ./cmd* -maxdepth 1 -mindepth 1 -type d -exec sh -c '"$(CURDIR)/scripts/docker-package.sh" kubecon {}' \;

.PHONY: build vet fmt fmtcheck errcheck package
