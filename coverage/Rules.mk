include mk/header.mk

$(d)/coverage_deps:
	-mkdir $(@D)/unitcover
	go get github.com/wadey/gocovmerge
	go get golang.org/x/tools/cmd/cover
.PHONY: $(d)/coverage_deps

UTESTS_$(d) := $(shell go list -f '{{if (len .TestGoFiles)}}{{.ImportPath}}{{end}}' $(go-flags-with-tags) ./... | grep -v go-ipfs/vendor | grep -v go-ipfs/Godeps)

UCOVER_$(d) := $(addsuffix .coverprofile,$(addprefix $(d)/unitcover/, $(subst /,_,$(UTESTS_$(d)))))

$(UCOVER_$(d)): $(d) $(d)/coverage_deps ALWAYS
	$(eval TMP_PKG := $(subst _,/,$(basename $(@F))))
	@go test $(go-flags-with-tags) $(GOTFLAGS) -covermode=atomic -coverprofile=$@ $(TMP_PKG)


$(d)/unit_tests.coverprofile: $(UCOVER_$(d))
	gocovmerge $^ > $@

TGTS_$(d) := $(d)/unit_tests.coverprofile

$(d)/sharness_tests.coverprofile: $(d)/coverage_deps

TGTS_$(d) += $(d)/sharness_tests.coverprofile

CLEAN += $(TGTS_$(d))
COVERAGE += $(TGTS_$(d))

include mk/footer.mk
