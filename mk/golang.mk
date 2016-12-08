# golang utilities
GOTAGS ?=
GOFLAGS ?=
GOTFLAGS ?=

go-pkg-name=$(shell go list ./$(1))
go-main-name=$(notdir $(call go-pkg-name,$(1)))$(?exe)
go-curr-pkg-tgt=$(d)/$(call go-main-name,$(d))

go-flags-with-tags=$(GOFLAGS)$(if $(GOTAGS), -tags $(call join-with,$(comma),$(GOTAGS)))

define go-build=
go build $(call go-flag-from-tags) -o "$@" "$(call go-pkg-name,$<)"
endef

.PHONY: test_go_short
test_go_short: GOTFLAGS += -test.short
test_go_short: test_go_expensive

.PHONY: test_go_race
test_go_race: GOTFLAGS += -race
test_go_race: test_go_expensive

TEST_GO := test_go_expensive
.PHONY: test_go_expensive
test_go_expensive:
	go test $(call go-flags-with-tags) $(GOTFLAGS) ./...

TEST_GO += test_go_fmt
.PHONY: test_go_fmt
test_go_fmt:
	bin/test-go-fmt

TEST += $(TEST_GO)


