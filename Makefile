.DEFAULT_GOAL := default

OPENCMD := open
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPENCMD = xdg-open
endif

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

up-s1:
	@docker-compose --file ./deployments/local/system-1/docker-compose.yml up --build -d
	$(OPENCMD) http://localhost:3000
	@docker logs -f system1_traffic-simulator_1

down-s1:
	@docker-compose --file ./deployments/local/system-1/docker-compose.yml down

up-s2:
	@docker-compose --file ./deployments/local/system-2/docker-compose.yml up --build -d
	@docker logs -f system2_wait-for-services_1
	@docker logs -f system2_traffic-simulator_1

down-s2:
	@docker-compose --file ./deployments/local/system-2/docker-compose.yml down

.PHONY: build vet fmt fmtcheck errcheck package up-s1
