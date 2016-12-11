
BINS = bin/random bin/multihash bin/ipfs bin/pollEndpoint bin/iptb bin/go-sleep
BINS += bin/go-timeout

export PATH := ../bin:${PATH}

# User might want to override those on the command line
GOFLAGS =

all: deps

global-deps:
	$(MAKE) -C .. deps

deps: global-deps bins

clean:
	$(MAKE) -C sharness clean
	rm -rf $(BINS)

bins: $(BINS)

# Tests

test: test_expensive

test_expensive: verify_gofmt
	$(MAKE) -C sharness TEST_EXPENSIVE=1
	$(MAKE) -C 3nodetest
	$(MAKE) -C dependencies

test_cheap: verify_gofmt
	$(MAKE) -C sharness
	$(MAKE) -C 3nodetest

test_race: verify_gofmt
	$(MAKE) -C sharness GOFLAGS=-race TEST_EXPENSIVE=1
	$(MAKE) -C 3nodetest GOFLAGS=-race
	$(MAKE) -C dependencies GOFLAGS=-race

# Misc

coverage: coverage_sharness

coverage_sharness:
	./sharness_test_coverage_helper.sh

IPFS-BUILD-OPTIONS: FORCE
	@bin/checkflags '$@' '$(GOFLAGS)' '*** new Go flags ***'

verify_gofmt:
	bin/verify-go-fmt.sh

.PHONY: all clean FORCE
