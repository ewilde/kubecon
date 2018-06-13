.DEFAULT_GOAL := default

OPENCMD := open
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPENCMD = xdg-open
endif

TEST?=$$(go list ./... |grep -v 'vendor')
GOIMPORT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: goimportscheck errcheck vet
	@find ./cmd/* -maxdepth 1 -type d -exec go install "{}" \;

goimports:
	goimports -w $(GOIMPORT_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

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
	@docker logs -f system-1_traffic-simulator_1

down-s1:
	@docker-compose --file ./deployments/local/system-1/docker-compose.yml down

up-s2:
	@docker-compose --file ./deployments/local/system-2/docker-compose.yml up --build -d
	@docker logs -f system-2_traffic-simulator_1

down-s2:
	@docker-compose --file ./deployments/local/system-2/docker-compose.yml down

.PHONY: build vet goimports goimportscheck errcheck package up-s1
