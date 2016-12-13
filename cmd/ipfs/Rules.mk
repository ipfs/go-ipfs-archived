include mk/header.mk
TGTS_$(d) := $(call go-curr-pkg-tgt)

TGT_BIN += $(TGTS_$(d))
CLEAN += $(TGTS_$(d))

PATH := $(realpath $(d)):$(PATH)

$(TGTS_$(d)): $(d) $$(DEPS_GO) ALWAYS # uses second expansion to collect all $(DEPS_GO)
	$(go-build)

include mk/footer.mk
