include mk/header.mk
TGTS_$(d) := $(call go-curr-pkg-tgt)

TGT_BIN += $(TGTS_$(d))
CLEAN += $(TGTS_$(d))

$(TGTS_$(d)): $(d) $$(GODEPS) ALWAYS # uses second expansion to collect all $(GODEPS)
	$(go-build)

include mk/footer.mk
