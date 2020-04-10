ALL_SRC := $(shell find . -name '*.go' \
								-not -path './vendor/*' \
								-not -path '*/gen-go/*' \
								-type f | sort)
ALL_PKGS := $(shell go list $(sort $(dir $(ALL_SRC))))

GOTEST_OPT?=-v -race -timeout 30s
GOTEST_OPT_WITH_COVERAGE = $(GOTEST_OPT) -coverprofile=coverage.txt -covermode=atomic
GOTEST=go test


.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf *.sql *.log test.db *coverage.out coverage.all integrations/*.sql

.PHONY: coverage
coverage:
	@hash gocovmerge > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/wadey/gocovmerge; \
	fi
	gocovmerge $(shell find . -type f -name "coverage.out") > coverage.all;\

.PHONY: travis-ci
travis-ci: test

.PHONY: test
test:
	$(GOTEST) $(GOTEST_OPT) $(ALL_PKGS)