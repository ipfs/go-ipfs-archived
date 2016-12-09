# golang utilities
GO_MIN_VERSION = 1.7

# pre-definitions
GOTAGS ?=
GOFLAGS ?=
GOTFLAGS ?=

GOFLAGS += -i # install packages to $GOPATH/pkg/ for major build speed-up

GODEPS :=

go-pkg-name=$(shell go list ./$(1))
go-main-name=$(notdir $(call go-pkg-name,$(1)))$(?exe)
go-curr-pkg-tgt=$(d)/$(call go-main-name,$(d))

go-flags-with-tags=$(GOFLAGS)$(if $(GOTAGS), -tags $(call join-with,$(comma),$(GOTAGS)))

define go-build=
go build $(go-flags-with-tags) -o "$@" "$(call go-pkg-name,$<)"
endef

test_go_short: GOTFLAGS += -test.short
test_go_short: test_go_expensive
.PHONY: test_go_short

test_go_race: GOTFLAGS += -race
test_go_race: test_go_expensive
.PHONY: test_go_race

test_go_expensive:
	go test $(go-flags-with-tags) $(GOTFLAGS) ./...
.PHONY: test_go_expensive
TEST_GO += test_go_expensive

test_go_fmt:
	bin/test-go-fmt
.PHONY: test_go_fmt
TEST_GO += test_go_fmt

test_go: $(TEST_GO)

check_go_version: 
	@bin/check_go_version $(IPFS_MIN_GO_VERSION)
.PHONY: check_go_version
CHECK_GO += check_go_version


TEST += $(TEST_GO)
